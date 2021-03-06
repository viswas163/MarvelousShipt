package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

var (
	// ComicsByCharacter : Map to store comics IDs by character name. map[char_name][comic_id] => map[string][[]int]
	ComicsByCharacter   sync.Map
	comicsByCharCounter = 0

	// ComicsByID : Map to get comics by their IDs. map[comic_id][comic] => map[int][models.Comic]
	ComicsByID        sync.Map
	comicsByIDCounter = 0
)

// Comic represents a Marvel comic.
type Comic struct {
	ID                 int            `json:"id,omitempty"`
	DigitalID          int            `json:"digitalId,omitempty"`
	Title              string         `json:"title,omitempty"`
	IssueNumber        int            `json:"issueNumber,omitempty"`
	VariantDescription string         `json:"variantDescription,omitempty"`
	Description        string         `json:"description,omitempty"`
	Format             string         `json:"format,omitempty"`
	PageCount          int            `json:"pageCount,omitempty"`
	ResourceURI        string         `json:"resourceURI,omitempty"`
	Variants           []ComicSummary `json:"variants,omitempty"`
	Collections        []ComicSummary `json:"collections,omitempty"`
	CollectedIssues    []ComicSummary `json:"collectedIssues,omitempty"`
	Characters         CharacterList  `json:"characters,omitempty"`
}

// ComicsDataWrapper provides character wrapper information returned by the API.
type ComicsDataWrapper struct {
	DataWrapper
	Data ComicsDataContainer `json:"data,omitempty"`
}

// ComicsDataContainer provides character container information returned by the API.
type ComicsDataContainer struct {
	DataContainer
	Results []Comic `json:"results,omitempty"`
}

// ComicList provides comics related to the parent entity.
type ComicList struct {
	List
	Items []ComicSummary `json:"items,omitempty"`
}

// ComicSummary provides the summary for a comic related to the parent entity.
type ComicSummary struct {
	Summary
}

// GetAllComics : Gets all comics from response
func GetAllComics(comicsJSON []byte) (ComicsDataWrapper, error) {
	var allComicsWrapper ComicsDataWrapper
	err := json.Unmarshal(comicsJSON, &allComicsWrapper)
	if err != nil {
		fmt.Println("Error unmarshaling all comics")
		return ComicsDataWrapper{}, err
	}
	return allComicsWrapper, nil
}

func checkComicsExist(cName string, allComs []Comic) bool {
	firstCID := allComs[0].ID
	lastCID := allComs[len(allComs)-1].ID

	secondOK, lastOK := false, false
	comics, firstOK := ComicsByCharacter.Load(cName)
	if firstOK {
		coms := reflect.ValueOf(comics)
		secondOK = coms.Index(0).Interface().(int) == firstCID
		lastOK = coms.Index(coms.Len()-1).Interface().(int) == lastCID
	}

	return firstOK && secondOK && lastOK
}

// SetAllComicsByCharName : Sets all comics to map by char name if not already existing
func SetAllComicsByCharName(cName string, allComics []Comic) error {
	if checkComicsExist(cName, allComics) {
		return nil
	}
	var comicIDsOfChar []int
	for _, comic := range allComics {
		comicIDsOfChar = append(comicIDsOfChar, comic.ID)
		ComicsByID.LoadOrStore(comic.ID, comic)
		comicsByIDCounter++
	}
	ComicsByCharacter.LoadOrStore(cName, comicIDsOfChar)
	comicsByCharCounter++
	fmt.Printf("Added %d new comics!", comicsByCharCounter)
	return nil
}
