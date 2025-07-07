package rss

import (
	"io"
	"net/http"
	"news-aggregator/pkg/storage"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		handler func(url string) (resp *http.Response, err error)
		url     string
	}

	testHandler := func(url string) (resp *http.Response, err error) {
		return &http.Response{
			Body: io.NopCloser(strings.NewReader(`<rss xmlns:dc="http://purl.org/dc/elements/1.1/" version="2.0">
<channel>
<title>
<![CDATA[ Все статьи подряд / Go / Хабр ]]>
</title>
<link>https://habr.com/ru/hubs/go/articles/</link>
<description>
<![CDATA[ Go – компилируемый, многопоточный язык программирования ]]>
</description>
<language>ru</language>
<managingEditor>editor@habr.com</managingEditor>
<generator>habr.com</generator>
<pubDate>Thu, 03 Jul 2025 15:02:00 GMT</pubDate>
<image>
...
</image>
<item>
<title>
<![CDATA[ FastCGo: как мы ускорили вызов C-кода в Go в 16,5 раза ]]>
</title>
<link>https://habr.com/ru/companies/flant/articles/923912</link>
<description>
<![CDATA[ <p>В Deckhouse Prom++ мы переписали ядро хранения и обработки горячих данных на C++, при этом вся оркестрация и периферия остались в Prometheus на Go, что позволило сохранить полную совместимость с Prometheus. Для частых вызовов кода C++ мы использовали механизм CGo, однако первые тесты показали, что производительность CPU практически не улучшилась из-за его медлительности. В итоге мы переписали CGo, создав собственный механизм вызова.</p><p>В статье разберём, что такое CGo и почему он такой медленный, сделаем простейший собственный механизм CGo-вызова и доведём этот механизм до полноценного решения.</p> <a href="https://habr.com/ru/articles/923912/?utm_campaign=923912&amp;utm_source=habrahabr&amp;utm_medium=rss#habracut">Читать далее</a> ]]>
</description>
<pubDate>Thu, 03 Jul 2025 05:58:40 GMT</pubDate>
</item>
</channel>
</rss>`)),
		}, nil
	}

	tests := []struct {
		name    string
		args    args
		want    []storage.Post
		wantErr bool
	}{
		{
			name: "Parse",
			args: args{
				handler: testHandler,
			},
			want: []storage.Post{
				{
					ID:      0,
					Title:   "\n FastCGo: как мы ускорили вызов C-кода в Go в 16,5 раза \n",
					Content: "\n В Deckhouse Prom++ мы переписали ядро хранения и обработки горячих данных на C++, при этом вся оркестрация и периферия остались в Prometheus на Go, что позволило сохранить полную совместимость с Prometheus. Для частых вызовов кода C++ мы использовали механизм CGo, однако первые тесты показали, что производительность CPU практически не улучшилась из-за его медлительности. В итоге мы переписали CGo, создав собственный механизм вызова.В статье разберём, что такое CGo и почему он такой медленный, сделаем простейший собственный механизм CGo-вызова и доведём этот механизм до полноценного решения. Читать далее \n",
					PubTime: 1751522320,
					Link:    "https://habr.com/ru/companies/flant/articles/923912",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.handler, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
