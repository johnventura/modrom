package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDiffsToJson(t *testing.T) {
	originalStr := []byte(
		`The second Conference at Basra, though many of its prominent figures had already played leading rôles in the earlier gathering in 1965, was very different in scale, scope and spirit from that assembly. It was an older gathering. The average age, says Amen Rihani, was a full ten years higher. Young men were still coming into the Fellowship abundantly, but there had also been accessions—and not always very helpful accessions—of older men who had been radical and revolutionary leaders in the war period. Their frame of experience had shaped them for irresponsible resistance. Their mental disposition was often obstructively critical and insubordinate. Many had had no sort of technical experience. They were disposed to throw an anarchistic flavour over schools and propaganda.

Moreover, the great scheme of the Modern State had now lost something of its first compelling freshness. The "young men of '65" had had ten years of responsible administrative work. They had been in contact with urgent detail for most of that period. They had had to modify De Windt's generalizations in many particulars, and the large splendour of the whole project no longer had the same dominating power over their minds. They had lost something of the professional esprit de corps, the close intimate confidence with each other, with which they had originally embarked upon the great adventure of the Modern State. Many had married women of the older social tradition and formed new systems of gratification and friendship. They had ceased to be enthusiastic young men and they had become men of the world. The consequent loss of a sure touch upon primary issues was particularly evident in the opening sessions.`)

	altStr := []byte(
		`The second Conference at Basra, though many of its prominent figures had already played leading rôles in the earlier gathering in 1965, was very different in scale, scope and spirit from that assembly. It was an older gathering. The average age, says Amen Rihani, was a full ten years higher. Young men were still coming into the Fellowship abundantly, but there had also been accessions—and not always very hurtful accessions—of older men who had been radical and revolutionary leaders in the war period. Their frame of experience had shaped them for irresponsible resistance. Their mental disposition was often obstructively critical and insubordinate. Many had had no sort of technical experience. They were disposed to throw an anarchistic flavour over schools and propaganda.

Moreover, the great scheme of the Modern State had now lost something of its first compelling freshness. The "young men of '65" had had ten years of responsible administrative work. They had been in contact with urgent detail for most of that period. They had had to modify De Windt's generalizations in many particulars, and the large splendour of the whole project no longer had the same dominating power over their minds. They had lost something of the professional esprit de corps, the close intimate confidence with each other, with which they had originally embarked upon the great adventure of the Modern State. Many had married women of the older social tradition and formed new systems of gratification and friendship. They had ceased to be enthusiastic young men and they had become men of the world. The consequent loss of a sure touch upon primary issues was particularly evident in the closing sessions.`)

	resultWeWant := `{"patch":[{"start":412,"offset":0,"newbytes":"757274","shasum":"843c295955090eb4b25b449f8b03b9ab88027d83ca888209eab597fc5149a67c","comment":""},{"start":1639,"offset":47,"newbytes":"636c6f73","shasum":"d7ab7c3249fca6159a6d786d7e9a09de83684378efca68ebc1374e01fbd499d5","comment":""}]}`

	reportStr, err := DiffsToJson(originalStr, altStr, int8(3))
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	assert.Equal(t, reportStr, resultWeWant, "Patch file generation worked")
}

func TestPatchByteArray(t *testing.T) {

	TestPatchFile := `
{"patch": [{"start": 940, "offset": 0, "newbytes": "706c6179", "shasum": "", "comment": "Chunk 1 changed"},{"start": 949, "offset": 0, "newbytes": "73", "shasum": "d5c4ad8394b778304665f8593d95a86ef144a5d40b0a0111b5099d11cad864e6", "comment": ""},{"start": 951, "end": 961, "offset": 0, "newbytes": "41746172692032363030", "shasum": "4d8277b2ab827426a81d80adebb954ea54b635ef9b2dbe8871ffc752787ea249", "comment": ""}]}`

	TextOriginal :=
		`In a castle of Westphalia, belonging to the Baron of
Thunder-ten-Tronckh, lived a youth, whom nature had endowed with the
most gentle manners. His countenance was a true picture of his soul. He
combined a true judgment with simplicity of spirit, which was the
reason, I apprehend, of his being called Candide. The old servants of
the family suspected him to have been the son of the Baron's sister, by
a good, honest gentleman of the neighborhood, whom that young lady would
never marry because he had been able to prove only seventy-one
quarterings, the rest of his genealogical tree having been lost through
the injuries of time.

The Baron was one of the most powerful lords in Westphalia, for his
castle had not only a gate, but windows. His great hall, even, was hung
with tapestry. All the dogs of his farm-yards formed a pack of hounds at
need; his grooms were his huntsmen; and the curate of the village was
his grand almoner. They called him "My Lord," and laughed at all his
stories.

The Baron's lady weighed about three hundred and fifty pounds, and was
therefore a person of great consideration, and she did the honours of
the house with a dignity that commanded still greater respect. Her
daughter Cunegonde was seventeen years of age, fresh-coloured, comely,
plump, and desirable. The Baron's son seemed to be in every respect
worthy of his father. The Preceptor Pangloss[1] was the oracle of the
family, and little Candide heard his lessons with all the good faith of
his age and character.

Pangloss was professor of metaphysico-theologico-cosmolo-nigology. He
proved admirably that there is no effect without a cause, and that, in
this best of all possible worlds, the Baron's castle was the most
magnificent of castles, and his lady the best of all possible
Baronesses.
`
	TextWeWant :=
		`In a castle of Westphalia, belonging to the Baron of
Thunder-ten-Tronckh, lived a youth, whom nature had endowed with the
most gentle manners. His countenance was a true picture of his soul. He
combined a true judgment with simplicity of spirit, which was the
reason, I apprehend, of his being called Candide. The old servants of
the family suspected him to have been the son of the Baron's sister, by
a good, honest gentleman of the neighborhood, whom that young lady would
never marry because he had been able to prove only seventy-one
quarterings, the rest of his genealogical tree having been lost through
the injuries of time.

The Baron was one of the most powerful lords in Westphalia, for his
castle had not only a gate, but windows. His great hall, even, was hung
with tapestry. All the dogs of his farm-yards formed a pack of hounds at
need; his grooms were his huntsmen; and the curate of the village was
his grand almoner. They played his Atari 2600 and laughed at all his
stories.

The Baron's lady weighed about three hundred and fifty pounds, and was
therefore a person of great consideration, and she did the honours of
the house with a dignity that commanded still greater respect. Her
daughter Cunegonde was seventeen years of age, fresh-coloured, comely,
plump, and desirable. The Baron's son seemed to be in every respect
worthy of his father. The Preceptor Pangloss[1] was the oracle of the
family, and little Candide heard his lessons with all the good faith of
his age and character.

Pangloss was professor of metaphysico-theologico-cosmolo-nigology. He
proved admirably that there is no effect without a cause, and that, in
this best of all possible worlds, the Baron's castle was the most
magnificent of castles, and his lady the best of all possible
Baronesses.
`
	patchJson, err := ReadPatchJson(TestPatchFile)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	testAlt, err := PatchByteArray([]byte(TextOriginal), patchJson)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	assert.Equal(t, string(testAlt), TextWeWant, "Patch output worked")
}
