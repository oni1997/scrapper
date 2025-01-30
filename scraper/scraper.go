package scraper

import (
    "context"
    "fmt"
    "log"
    "strings"
    "time"

    "github.com/chromedp/chromedp"
    "scrapper/models"
)

type Scraper struct {
    baseURL string
    timeout time.Duration
    retries int
}

func NewScraper() *Scraper {
    return &Scraper{
        baseURL: "https://www.pnp.co.za",
        timeout: 60 * time.Second,
        retries: 3,
    }
}

func (s *Scraper) waitForProducts() chromedp.Action {
    return chromedp.ActionFunc(func(ctx context.Context) error {
        if err := chromedp.WaitReady("body", chromedp.ByQuery).Do(ctx); err != nil {
            return err
        }

        var pageSource string
        if err := chromedp.OuterHTML("html", &pageSource).Do(ctx); err != nil {
            return err
        }

        time.Sleep(5 * time.Second)
        return nil
    })
}

func (s *Scraper) ScrapeProducts(searchTerm string) (*models.ProductResponse, error) {
    var lastErr error
    
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", true),
        chromedp.Flag("disable-gpu", true),
        chromedp.Flag("no-sandbox", true),
        chromedp.Flag("disable-dev-shm-usage", true),
        chromedp.Flag("disable-software-rasterizer", true),
        chromedp.Flag("in-process-gpu", true),
        chromedp.Flag("disable-setuid-sandbox", true),
    )

    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()

    for attempt := 0; attempt < s.retries; attempt++ {
        if attempt > 0 {
            log.Printf("Retry attempt %d/%d", attempt+1, s.retries)
            time.Sleep(time.Duration(attempt) * 2 * time.Second)
        }

        ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
        ctx, cancel = context.WithTimeout(ctx, s.timeout)
        defer cancel()

        var products []models.Product
        searchURL := fmt.Sprintf("%s/search/%s", s.baseURL, searchTerm)

        err := chromedp.Run(ctx,
            chromedp.Navigate(searchURL),
            chromedp.ActionFunc(func(ctx context.Context) error {
                log.Println("Page navigation complete, waiting for content...")
                return nil
            }),
            chromedp.Sleep(5*time.Second),
            s.waitForProducts(),
            chromedp.ActionFunc(func(ctx context.Context) error {
                log.Println("Attempting to extract products...")
                return nil
            }),
            chromedp.Evaluate(`
                (() => {
                    const products = document.getElementsByClassName('product-grid-item');
                    console.log('Found products:', products.length);
                    return Array.from(products).map(item => {
                        const img = item.querySelector('.product-grid-item__image-container.product-action');
                        const promo = item.querySelector('.product-grid-item__promotion-container a');
                        const price = item.querySelector('.price');
                        
                        return {
                            id: item.getAttribute('data-cnstrc-item-id') || '',
                            name: item.getAttribute('data-cnstrc-item-name') || '',
                            price: price ? price.textContent.trim() : '',
                            imageUrl: img ? img.getAttribute('src') : '',
                            promotion: promo ? promo.textContent.trim() : null
                        };
                    });
                })()
            `, &products),
        )

        if err != nil {
            lastErr = fmt.Errorf("attempt %d failed: %v", attempt+1, err)
            log.Printf("Scraping error: %v", lastErr)
            continue
        }

        for i := range products {
            if products[i].Price != "" {
                products[i].Price = strings.TrimSpace(products[i].Price)
                if !strings.HasPrefix(products[i].Price, "R") {
                    products[i].Price = "R" + products[i].Price
                }
            }
        }

        if len(products) == 0 {
            lastErr = fmt.Errorf("no products found on attempt %d", attempt+1)
            log.Println(lastErr)
            continue
        }

        log.Printf("Successfully scraped %d products", len(products))
        return &models.ProductResponse{
            Success:    true,
            Products:   products,
            Total:     len(products),
            SearchTerm: searchTerm,
        }, nil
    }

    return &models.ProductResponse{
        Success:    false,
        Message:    fmt.Sprintf("Failed after %d attempts. Last error: %v", s.retries, lastErr),
        SearchTerm: searchTerm,
    }, lastErr
}