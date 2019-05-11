package diplomat

import (
	"github.com/tony84727/diplomat/pkg/data"
	"github.com/tony84727/diplomat/pkg/emit"
	"github.com/tony84727/diplomat/pkg/log"
	"github.com/tony84727/diplomat/pkg/selector"
	"strings"
	"sync"
)

type Synthesizer struct {
	data.Translation
	output Output
	emitterRegistry emit.Registry
	logger          log.Logger
}

func NewSynthesizer(output Output, translation data.Translation, emitterRegistry emit.Registry, logger log.Logger) *Synthesizer {
	return &Synthesizer{translation, output, emitterRegistry, log.MaybeLogger(logger)}
}

func (s Synthesizer) Output(output data.Output) error {
	selectors := output.GetSelectors()
	selectorInstance := make([]selector.Selector, len(selectors))
	for i, s := range selectors {
		selectorInstance[i] = selector.NewPrefixSelector(strings.Split(string(s), ".")...)
	}
	selected := data.NewSelectedTranslation(s, selector.NewCombinedSelector(selectorInstance...))
	templates := output.GetTemplates()
	errChan := make(chan error)
	var wg sync.WaitGroup
	for _, t := range templates {
		if i := s.emitterRegistry.Get(t.GetType()); i != nil {
			wg.Add(1)
			go func(t data.Template) {
				defer wg.Done()
				s.logger.Info("[Emitting] %s [%s]", t.GetOptions().GetFilename(), t.GetType())
				output, err := i.Emit(selected, t.GetOptions())
				if err != nil {
					errChan <- err
					return
				}
				if err := s.output.WriteFile(t.GetOptions().GetFilename(), output); err != nil {
					errChan <- err
				}
			}(t)
		}
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	defer func() {
		// dump error
		go func() {
			for range errChan {
			}
		}()
		wg.Wait()
		close(errChan)
	}()
	select {
	case err := <-errChan:
		s.logger.Error(err.Error())
		return err
	case <-done:
		return nil
	}
}
