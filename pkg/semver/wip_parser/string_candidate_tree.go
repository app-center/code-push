package wip_parser

type stringCandidateNode struct {
	char rune

	left  *stringCandidateNode
	right *stringCandidateNode
}

type stringCandidateTree struct {
	head *stringCandidateNode
}

func (t stringCandidateTree) Walk(c rune) (subTree stringCandidateTree, hit bool) {
	emptyTree := stringCandidateTree{}

	walkNode := t.head

	for walkNode != nil {
		if walkNode.char == c {
			return stringCandidateTree{head: walkNode.right}, true
		}
		walkNode = walkNode.left
	}

	return emptyTree, false
}

func newStringCandidateTree(candies []string) stringCandidateTree {
	tree := stringCandidateTree{}

	root := &stringCandidateNode{char: ' '}

	for _, candi := range candies {
		if candi == "" {
			continue
		}

		parentNode := root
		for _, char := range candi {
			parentNode = pushCharNode(parentNode, char)
		}
	}

	tree.head = root.right

	return tree
}

func pushCharNode(parentNode *stringCandidateNode, char rune) *stringCandidateNode {
	charNode := &stringCandidateNode{char: char}

	if parentNode == nil {
		return charNode
	}

	if parentNode.right == nil {
		parentNode.right = charNode
		return charNode
	} else if char > parentNode.right.char {
		charNode.right = parentNode.right.right
		charNode.left = parentNode.right
		parentNode.right = charNode
		return charNode
	} else if char < parentNode.right.char {
		prevNode := parentNode.right

		for true {
			nextNode := prevNode.left

			if nextNode == nil {
				prevNode.left = charNode
				return charNode
			}

			if nextNode.char > char {
				prevNode = nextNode
				continue
			} else if nextNode.char == char {
				return nextNode
			} else if nextNode.char < char {
				prevNode.left = charNode
				charNode.left = nextNode
				return charNode
			}
		}
	} else if char == parentNode.right.char {
		return parentNode.right
	}

	return charNode
}
