# CloudWatchExporterSaver
CloudWatchExporterSaver aims to reduce the request frequency of the CloudWatch exporter to minimize operational costs while maintaining effective monitoring. 

## Usage

```bash
go run main.go
```

## Configuration

```bash
export CLOUDWATCH_EXPORTER_URL=http://localhost:9106/metrics
```

## Metrics

```bash
curl http://localhost:9107/metrics
```