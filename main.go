package main

import (
	"context"
	"fmt"
	"bufio"
	"os"
	"strings"
	"log"
	
	
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Reverse the array of versions so as to make it in descending order
func ReverseArray(s []*semver.Version) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    	s[i], s[j] = s[j], s[i]
	}
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	var maxPatch *semver.Version

	if (len(releases) == 0){
		return versionSlice
	}

	semver.Sort(releases)
	ReverseArray(releases)
	maxPatch = releases[0]
	
	for _, release := range releases {
		if i := (*release).Compare(*minVersion); i == 1 {
			if (*maxPatch).Major != (*release).Major || (*maxPatch).Minor != (*release).Minor {
				versionSlice = append(versionSlice, maxPatch)
				maxPatch = release
			}
		} else {
			break
		}
	}

	if (len(versionSlice) == 0){
		versionSlice = append(versionSlice, maxPatch)
	} else if (versionSlice[len(versionSlice) - 1] != maxPatch) {
		versionSlice = append(versionSlice, maxPatch)
	}
	
	return versionSlice
}

func TackleEachApplication(information string) {
	info := strings.Split(information, ",")
	respository := strings.Split(info[0], "/")

    // Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, respository[0], respository[1], opt)
	if err != nil {
		log.Fatal(err) // is this really a good way?
	}
	minVersion := semver.New(info[1])
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	versionSlice := LatestVersions(allReleases, minVersion)

	fmt.Printf("latest versions of %s/%s: %s\n", respository[0], respository[1], versionSlice)
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	filePath := os.Args[1]
	f, err := os.Open(filePath)
	check(err)

    fileScanner := bufio.NewScanner(f)
    // I'm not sure whether the first line of the file will be a redundant "repository,min_version" as in the example, 
    // or some actual data that needs to be processed. I assumed that there is no redundant line in this file. 
	for fileScanner.Scan() {
		TackleEachApplication(fileScanner.Text())
	}

	if err := fileScanner.Err(); err != nil {
	    log.Fatal(err)
	}
}