package wip_parser

import "fmt"

type iSegment interface {
	walk(c rune)(hit bool)
	report() error
}

type segmentConfig struct {
	walker   iWalker
	reporter segReporter
}

type switchSegmentCases func(c rune) iSegment
var switchSegmentDefaultCase = func(c rune) iSegment {
	return nil
}

type switchSegmentConfig struct {
	cases    switchSegmentCases
	reporter segReporter
}

type chainSegmentConfig struct {
	chains   []iSegment
	reporter segReporter
}

type immediateSegmentConfig struct {
	reporter segReporter
}

type segment struct {
	walker   iWalker
	reporter segReporter

	chars []rune
}

func (s *segment) walk(c rune) (hit bool) {
	defer func() {
		if r := recover(); r != nil {
			s.reporter(fmt.Errorf("panic error: %v", r), nil)
			hit = false
		}
	}()

	if s.walker(c) {
		s.chars = append(s.chars, c)
		hit = true
	} else {
		hit = false
	}

	return
}

func (s *segment) report() error {
	if s.reporter == nil {
		return nil
	}

	return s.reporter(nil, s.chars)
}

type switchSegment struct {
	reporter segReporter
	cases    switchSegmentCases

	used iSegment
}

func (s *switchSegment) walk(c rune) (hit bool) {
	if s.used != nil {
		return s.used.walk(c)
	}

	seg := s.cases(c)

	if seg != nil {
		s.used = seg
		hit = true
	} else {
		hit = false
	}

	return
}

func (s *switchSegment) report() error {
	if s.used == nil {
		return s.reporter(nil, nil)
	}

	return s.used.report()
}

type chainSegment struct {
	queueSegments []iSegment
	walkedSegments []iSegment
	walkingSegment iSegment

	reporter   segReporter
	walkedChar []rune
}

func (cs *chainSegment) walk(c rune) (hit bool) {
	cs.walkedChar = append(cs.walkedChar, c)

	for true {
		if cs.walkingSegment == nil {
			if len(cs.queueSegments) == 0 {
				return false
			} else {
				cs.walkingSegment = cs.queueSegments[0]
				cs.queueSegments = cs.queueSegments[:len(cs.queueSegments)-1]
			}
		}

		walkingSegment := cs.queueSegments[0]
		if walkingSegment.walk(c) {
			return true
		} else {
			cs.walkedSegments = append(cs.walkedSegments, walkingSegment)

			if len(cs.queueSegments) > 1 {
				cs.queueSegments = cs.queueSegments[1:]
				if !cs.queueSegments[0].walk(c) {
					cs.walkedSegments = append(cs.walkedSegments, cs.queueSegments[0])
					if len(cs.queueSegments) > 1 {
						cs.queueSegments = cs.queueSegments[1:]
					}
					return false
				} else {
					return true
				}
			} else {
				return false
			}
		}
	}

	return false
}

func (cs *chainSegment) report() error {
	for i := range cs.walkedSegments {
		if subErr := cs.walkedSegments[i].report(); subErr != nil {
			return subErr
		}
	}

	if len(cs.queueSegments) > 0 {
		return cs.reporter(fmt.Errorf("segment not preciesly matched, walked char: %v", cs.walkedChar), nil)
	}

	return nil
}

type immediateSegment struct {
	reporter segReporter
}

func (is *immediateSegment) walk(char rune) bool {
	return false
}

func (is *immediateSegment) report() error {
	return is.reporter(nil, nil)
}

func newSegment(config segmentConfig) *segment {
	if config.walker == nil {
		config.walker = noopWalker
	}

	if config.reporter == nil {
		config.reporter = noopSegReporter
	}

	return &segment{
		walker:   config.walker,
		reporter: config.reporter,
		chars:    nil,
	}
}

func newSwitchSegment(config switchSegmentConfig) *switchSegment {
	if config.reporter == nil {
		config.reporter = noopFailedSegReporter
	}

	if config.cases == nil {
		config.cases = switchSegmentDefaultCase
	}

	return &switchSegment{
		cases:   config.cases,
		reporter: config.reporter,
	}
}

func newChainSegment(config chainSegmentConfig) *chainSegment {
	if config.reporter == nil {
		config.reporter = noopSegReporter
	}

	return &chainSegment{
		queueSegments:  config.chains,
		reporter:       config.reporter,
	}
}

func newImmediateSegment(config immediateSegmentConfig) *immediateSegment {
	if config.reporter == nil {
		config.reporter = noopSegReporter
	}

	return &immediateSegment{reporter: config.reporter}
}
