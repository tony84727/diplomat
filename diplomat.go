package diplomat

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/siongui/gojianfan"
)

type Diplomat struct {
	outline Outline
}

func NewDiplomatForFile(path string) (*Diplomat, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var outline Outline
	err = yaml.Unmarshal(data, &outline)
	if err != nil {
		return nil, err
	}
	dip := NewDiplomat(outline)
	return &dip, nil
}

func NewDiplomat(outline Outline) Diplomat {
	d := Diplomat{
		outline,
	}
	d.applyTransformers()
	return d
}

func (d Diplomat) applyTransformers() {
	if d.outline.Settings.Chinese != nil {
		d.applyChineseConvertor(
			d.outline.Settings.Chinese.Convert.Mode,
			d.outline.Settings.Chinese.Convert.From,
			d.outline.Settings.Chinese.Convert.To,
		)
	}
}

type ChineseConvertor struct {
	From          string
	To            string
	transformFunc func(string) string
}

func NewChineseConvertor(mode, from, to string) (*ChineseConvertor, error) {
	var transformFunc func(string) string
	switch mode {
	case "s2f":
		transformFunc = gojianfan.S2T
		break
	case "t2s":
		transformFunc = gojianfan.T2S
		break
	default:
		return nil, fmt.Errorf("chinese convertor: unknown mode %s", mode)
	}
	return &ChineseConvertor{
		transformFunc: transformFunc,
		From:          from,
		To:            to,
	}, nil
}

func (cc ChineseConvertor) Apply(t Translation) {
	if cc.appliable(t) {
		from, _ := t.Get(cc.From)
		t.Set(cc.To, cc.transformFunc(from))
	}
}

func (cc ChineseConvertor) appliable(t Translation) bool {
	_, fromExist := t.Get(cc.From)
	_, toExist := t.Get(cc.To)
	return fromExist && !toExist
}

func (d *Diplomat) applyChineseConvertor(mode, from, to string) error {
	convertor, err := NewChineseConvertor(mode, from, to)
	if err != nil {
		return err
	}
	for _, f := range d.outline.Fragments {
		for _, t := range f.Translations {
			convertor.Apply(t)
		}
	}
	return nil
}

func (d Diplomat) GetOutline() Outline {
	return d.outline
}
