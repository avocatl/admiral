package display

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var currencies = []map[string]interface{}{
	{"Symbol": "EUR", "Quote": 1},
	{"Symbol": "USD", "Quote": 1.22},
	{"Symbol": "MXN", "Quote": 24.45},
}

var currencyColMap = map[string]string{
	"Symbol": "The currency symbol",
	"Quote":  "The current exchange rate",
}

var currencyCol = []string{
	"Symbol", "Quote",
}

func TestDisplay_DefaultDisplayer_NoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDisplayable(ctrl)

	m.EXPECT().KV().Return(currencies)
	m.EXPECT().Cols().AnyTimes().Return(currencyCol)
	m.EXPECT().NoHeaders().Return(false)

	dsp := DefaultDisplayer(nil)
	got := dsp.Display(m, []string{})

	assert.Nil(t, got)
}

func TestDisplay_DefaultDisplayer_Content(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDisplayable(ctrl)

	m.EXPECT().KV().Return(currencies)
	m.EXPECT().Cols().AnyTimes().Return(currencyCol)
	m.EXPECT().NoHeaders().Return(false)

	want := `Symbol    Quote
EUR       1
USD       1.220000
MXN       24.450000
`

	b := bytes.NewBufferString("")

	dsp := DefaultDisplayer(b)
	_ = dsp.Display(m, []string{})

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, want, string(out))
}

func TestDisplay_Displayable_ColMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDisplayable(ctrl)

	m.EXPECT().ColMap().MaxTimes(1).Return(currencyColMap)

	assert.Contains(t, m.ColMap(), "Quote")
}

func TestDisplay_DefaultDisplayer_Filtered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDisplayable(ctrl)

	m.EXPECT().KV().Return(currencies)
	m.EXPECT().Cols().AnyTimes().Return(currencyCol)
	m.EXPECT().NoHeaders().Return(false)

	want := `Symbol
EUR
USD
MXN
`
	b := bytes.NewBufferString("")

	dsp := DefaultDisplayer(b)
	_ = dsp.Display(m, []string{"Symbol"})

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, want, string(out))

}

func TestDisplay_DefaultDisplayer_NoHeaders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDisplayable(ctrl)

	m.EXPECT().KV().Return(currencies)
	m.EXPECT().Cols().AnyTimes().Return(currencyCol)
	m.EXPECT().NoHeaders().Return(true)

	want := "EUR    1\nUSD    1.220000\nMXN    24.450000\n"

	b := bytes.NewBufferString("")

	dsp := DefaultDisplayer(b)
	_ = dsp.Display(m, []string{})

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, want, string(out))
}

func TestFilterColumns(t *testing.T) {
	cases := []struct {
		name  string
		given string
		def   []string
		want  []string
	}{
		{
			"filterable columns passed in lower",
			"id,name",
			[]string{"id", "name", "address", "surname"},
			[]string{"id", "name"},
		},
		{
			"no columns passed returns default",
			"",
			[]string{"id", "name", "address", "surname"},
			[]string{"id", "name", "address", "surname"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := FilterColumns(c.given, c.def)
			assert.Equal(t, got, c.want)
		})
	}
}
