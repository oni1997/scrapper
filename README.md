# PnP Product Scraper API

A robust Go-based web scraping API that fetches product information from Pick n Pay's website. This API provides endpoints to search for products and retrieve their details including prices, promotions, and images.

## Features

* Real-time web scraping using ChromeDP for JavaScript-rendered content
* Configurable retry mechanism for reliable data collection
* RESTful API endpoints with JSON responses
* Headless browser automation for efficient scraping
* Automatic price formatting and data normalization
* Comprehensive error handling and reporting

## Prerequisites

* Go 1.16 or higher
* Chrome/Chromium browser
* Git

## Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:oni1997/scrapper.git
   cd scrapper
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Configuration

The application uses the following default configuration:

* Base URL: `https://www.pnp.co.za`
* Request Timeout: 60 seconds
* Maximum Retries: 3
* Server Port: 8080

## API Reference

### Search Products

Searches for products based on the provided search term.

* **URL**: `/api/search`
* **Method**: `POST`
* **Content-Type**: `application/json`

**Request Body**:
```json
{
    "searchTerm": "string"
}
```

**Success Response**:
* **Code**: 200
* **Content**:
  ```json
  {
      "success": true,
      "message": "string",
      "products": [
          {
              "id": "string",
              "name": "string",
              "price": "string",
              "image_url": "string",
              "promotion": "string"
          }
      ],
      "total": 0,
      "search_term": "string"
  }
  ```

**Error Responses**:
* **Code**: 400 Bad Request
* **Code**: 405 Method Not Allowed
* **Code**: 500 Internal Server Error

## Usage Example

```bash
curl -X POST http://localhost:8080/api/search \
  -H "Content-Type: application/json" \
  -d '{"searchTerm": "milk"}'
```

## Technical Architecture

### Core Components

1. **Scraper Engine**
   * Manages ChromeDP instances for web scraping
   * Implements intelligent retry mechanisms
   * Handles DOM navigation and content extraction

2. **Data Models**
   * `Product`: Product information structure
   * `ProductResponse`: API response wrapper
   * `SearchRequest`: Search request structure

3. **API Server**
   * HTTP request handling
   * Route management
   * Response formatting
   * Error handling

### Browser Automation Configuration

The ChromeDP instance runs with the following settings:

* Headless mode enabled
* GPU disabled
* Sandbox disabled
* Custom user agent
* Window size: 1920x1080

## Development

Run the development server:

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

## Production Deployment

Consider the following when deploying to production:

* Configure appropriate timeout values
* Adjust retry mechanisms based on server capacity
* Implement rate limiting
* Add authentication layer
* Enable HTTPS
* Set up monitoring and logging
* Implement caching mechanisms

## Contributing

Contributions are welcome! Please feel free to submit pull requests.

## Author

Onesmus

## License

MIT License

---

For bug reports and feature requests, please open an issue on the GitHub repository.