# AWS S3 Uploader API


# API

The REST API to the example app is described below.

## Get list of files

### Request

`GET /{path}`

    curl -i -H 'Accept: application/json' http://127.0.0.1/{path}

### Response

  HTTP/1.1 200 OK
  Content-Type: application/json
  Date: Wed, 20 Nov 2019 01:53:41 GMT
  Content-Length: 51


```json
    {
      "files": ["<filePath>"]
    }
```

## Upload a file to S3

### Request

`POST /`

    curl -i -H 'Accept: application/json' http://127.0.0.1:3000 \
      -d '{"path": "./somePath", "bucket": "bucket-name", "region": "us-east-1"}'

### Response

  HTTP/1.1 200 OK
  Content-Type: application/json
  Date: Wed, 20 Nov 2019 01:52:28 GMT
  Content-Length: 23

```json
    {
      "message": [ "Updated" ]
    }
```

## Reference

- https://docs.aws.amazon.com/code-samples/latest/catalog/go-s3-s3_upload_directory.go.html
