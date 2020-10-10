# LinksChecker

![Go](https://github.com/Ythosa/linkschecker/workflows/Go/badge.svg)

## Description
LinkChecker is service for checking links. It can be used in large systems as a microservice that can be upped as a regular Docker container.

## Installation
  * Clone from Git: 
    ```bash
    $ git clone https://github.com/Ythosa/linkschecker
    ```
  * Transfer environment variables:
    * PORT;
    * LOG_LEVEL;
  * Start service:
    * As Docker container:
      ```bash
      $ docker-compose up --build
      ```
    * Using Makefile:
      ```bash
      $ make run
      ```
    * Using Go CLI:
      ```bash
      $ go run ./src/cmd/apiserver/main.go
      ```
  * It's all :)

## API Methods

<table border="0.2">
<!--    <caption>Таблица размеров обуви</caption> -->
   <tr>
    <th>Method</th>
    <th>Description</th>
    <th>Request</th>
    <th>Response</th>
   </tr>
   <tr>
     <td> /get_broken_links </td>
     <td> Returns all broken links from site </td>
     <td> `{"base_url": "some_url"}`
     <td> `{"broken_links": {"url": "error|null"}}`
   </tr>
   <tr>
     <td> /validate_link </td>
     <td> Validates link and returns `ok` and `error` </td>
     <td> `{"link": "some_url"}`
     <td> `{"ok":"true|false", "error":""}`
   </tr>
   <tr>
        <td> /validate_link </td>
        <td> Validates list of links and returns list of `url` and `error` for each link</td>
        <td> `{"links": ["some_url1", "some_url2", ...]}`
        <td> `[{"url": "...", "error": "..."}, ...]`
  </tr>
</table>
