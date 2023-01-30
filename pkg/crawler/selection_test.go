package crawler

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestSelection(t *testing.T) {
	t.Run("should do convert", func(t *testing.T) {
		d := DocumentFromBytes([]byte(`
<aside class="pager">
    <ul>
        <li><a href="?page=1">1</a></li>
        <li><a href="?page=2" class="active">2</a></li>
        <li><a href="?page=3">3</a></li>
    </ul>
</aside>
`))

		s := must(SelectionFromGraphQL([]byte(`
query {
	page @_dom_query(select: "aside.pager ul a.active") @_dom_get(idx: 0)
}

type Query {
	page: Int
}
`)))
		data, err := s.Extract(d.Selection)
		testingx.Expect(t, err, testingx.Be[error](nil))
		testingx.Expect(t, data, testingx.Equal[any](object{
			"page": 2,
		}))
	})

	t.Run("should do list", func(t *testing.T) {
		d := DocumentFromBytes([]byte(`
<main>
    <div class="list">
        <a class="list__item" href="/articles/1">
            <div class="title">
                一
            </div>
        </a>
        <a class="list__item" href="/articles/2">
            <div class="title">
                二
            </div>
        </a>
        <a class="list__item" href="/articles/3">
            <div class="title">
                三
            </div>
        </a>
    </div>
</main>
`))

		s := must(SelectionFromGraphQL([]byte(`
query {
	data @_dom_query(select: "a.list__item") {
		title @_dom_query(select: ".title")	
		id @_dom_attr(name: "href") @_string_replace(fromR: "/articles\/([^/]+)", to: "$1")
	}
}

type Query {
	data: [Item]
}

type Item {
	title: String
	id: String
}
`)))

		data, err := s.Extract(d.Selection)
		testingx.Expect(t, err, testingx.Be[error](nil))
		testingx.Expect(t, data, testingx.Equal[any](object{
			"data": list{
				object{
					"id":    "1",
					"title": "一",
				},
				object{
					"id":    "2",
					"title": "二",
				},
				object{
					"id":    "3",
					"title": "三",
				},
			},
		}))
	})

	t.Run("should do related query", func(t *testing.T) {
		d := DocumentFromBytes([]byte(`
<aside class="filters">
	<ul>
		<li><a href="#cat-a">分类一</a></li>
		<li><a href="#cat-b" class="active">分类二</a></li>
	</ul>
	<div>
		<div id="cat-a">
			<ul>
				<li><a href="">全部</a></li>
				<li><a href="?a=1">条件一</a></li>
				<li><a href="?a=2">条件二</a></li>
			</ul>
		</div>
		<div id="cat-b">
			<ul>
				<li><a href="">全部</a></li>
				<li><a href="?b=1">条件一</a></li>
				<li><a href="?b=2">条件二</a></li>
			</ul>
		</div>
	</div>
</aside>
`))

		s := must(SelectionFromGraphQL([]byte(`
query {
	data @_dom_query(select: "[id^=cat-]") {
		subgroups @_dom_query(select: "a") {
			name
			link @_dom_attr(name:"href")
		}
		name 
			@_dom_attr(name: "id") @_def(var: "id") 
			@_dom_closet(select: ".filters") @_dom_query(select: "a[href='#{id}']")
	}
}

type Query {
	data: [Group]
}

type Group {
	name: String
	subgroups: [SubGroup]
}

type SubGroup {
	name: String
	link: String
}
`)))

		data, err := s.Extract(d.Selection)
		testingx.Expect(t, err, testingx.Be[error](nil))

		testingx.Expect(t, data, testingx.Equal[any](object{
			"data": list{
				object{
					"name": "分类一",
					"subgroups": list{
						object{"name": "全部", "link": ""},
						object{"name": "条件一", "link": "?a=1"},
						object{"name": "条件二", "link": "?a=2"},
					},
				},
				object{
					"name": "分类二",
					"subgroups": list{
						object{"name": "全部", "link": ""},
						object{"name": "条件一", "link": "?b=1"},
						object{"name": "条件二", "link": "?b=2"},
					},
				},
			},
		}))
	})
}

func must[R any](input R, err error) R {
	if err != nil {
		panic(err)
	}
	return input
}

type object = map[string]any
type list = []any
