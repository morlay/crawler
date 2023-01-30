CRAWLER = go run ./cmd/crawler
CRAWLER_DEBUG = $(CRAWLER) -s=.tmp/source.gql


test:
	go test -v ./...

lint:
	goimports -w -l ./cmd
	goimports -w -l ./pkg


debug.list:
	$(CRAWLER_DEBUG) do getVideoList -p=cat=2 -p=page=2

debug.one:
	$(CRAWLER_DEBUG) do getVideoDetail -p=video_id=3777

#	http://0.0.0.0:7666?operation=getVideoList&cat=2
debug.serve:
	$(CRAWLER_DEBUG) serve
