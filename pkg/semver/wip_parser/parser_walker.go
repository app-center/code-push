package wip_parser

type iWalker func(c rune) (hit bool)

func noopWalker (c rune) (hit bool) {
	return false
}

func digitWalker(c rune) (hit bool) {
	return c >= '0' && c <= '9'
}

func dotWalker(c rune) bool {
	return c == '.'
}

func dashWalker(c rune) bool {
	return c == '-'
}

func eolWalker(c rune) bool {
	return c == 0x00
}

func charWalkerFactory(t rune) iWalker {
	return func(c rune) (hit bool) {
		return c == t
	}
}

func lengthLimitedWalkerFactory(length uint, walker iWalker) iWalker {
	remain := length

	return func(c rune) (hit bool) {
		if !walker(c) {
			return false
		}

		if remain <= 0 {
			return false
		} else {
			remain--
			return true
		}
	}
}

func digitWithLenWalkerFactory(length uint) iWalker {
	return lengthLimitedWalkerFactory(length, digitWalker)
}

func stringCandidateWalkerFactory(candies []string) iWalker {
	tree := newStringCandidateTree(candies)

	return func(c rune) (hit bool) {
		subTree, hit := tree.Walk(c)

		if hit {
			tree = subTree
			return true
		} else {
			return false
		}
	}
}