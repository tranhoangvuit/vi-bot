package vibot

func (vb *ViBot) SearchDictionary(wordSearch string) {
	vb.GetDefinition("http://www.thefreedictionary.com/", wordSearch)
}
