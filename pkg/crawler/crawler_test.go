package crawler

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

var code = `
type Query{
    from(origin: String!, name: String): Operations
}

type Operations {
    getVideoList(path: String!, parameters: [Parameter]): VideoList
}

enum Location {
    Path
    Query
}

input Option {
	v: String
	n: String
}

input Parameter {
    name: String!
    in: Location
    default: String
    enum: [Option]
}

type VideoList {
    data: [VideoSummary]
    page: Int
}

type VideoSummary  {
    video_id: String
    title: String
    cover: String
}

query {
	from(origin: "https://test.io", name: "x") {
		getVideoList(
			path: "/vodshow/{cat}-{area}-{by}-{class}-----{page}---{year}/"
			parameters: [
                {
                    name:"cat",
                    default:"1",
                    enum: [
                        {v: "1", n:"Movie"} 
                        {v: "4", n:"Anime"} 
                    ],
                },
                {name:"area"},
                {name:"by"},
                {name:"class"},
                {name:"year"},
                {name:"page"},
            ]
		) {
			data @_dom_query(select: "a.myui-vodlist__thumb") {
                title 
                video_id @_dom_attr(name: "href") @_string_replace(fromR: "/voddetail/([^/]+)/", to: "$1" )
                cover @_dom_attr(name: "data-original")
            }
            page @_dom_query(select: ".myui-page a.btn-warm") @_dom_get(idx: 0)
		}
	}
}
`

func Test(t *testing.T) {
	c, err := NewCrawler([]byte(code))
	testingx.Expect(t, err, testingx.Be[error](nil))

	s := c.Source()

	testingx.Expect(t, s.Name, testingx.Be("x"))
	testingx.Expect(t, s.Origin, testingx.Be("https://test.io"))
	testingx.Expect(t, len(s.Operations), testingx.Be(1))

	t.Run("should get correct path", func(t *testing.T) {
		op := s.Operations["getVideoList"]

		u, err := op.RequestURI(s, map[string]string{
			"cat": "2",
		})
		testingx.Expect(t, err, testingx.Be[error](nil))
		testingx.Expect(t, u.String(), testingx.Be("https://test.io/vodshow/2-----------/"))
	})
}
