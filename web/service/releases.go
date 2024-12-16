package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"x-ui/util/common"
	"x-ui/xray"
)

var (
	releaseLock    sync.Mutex
	releaseService *ReleaseService
)

const (
	maxTrojanVer    = "v1.8"
	minSupportedVer = "v1.4.2"
)

func GetReleaseService() *ReleaseService {
	releaseLock.Lock()
	defer releaseLock.Unlock()
	if releaseService == nil {
		releaseService = &ReleaseService{
			cacheLock: sync.Mutex{},
			cache:     make(map[string]Release),
		}
		err := releaseService.loadRelease(xray.GetReleasePath())
		if err != nil {
			releaseService = nil
		}
	}
	return releaseService
}

type ReleaseService struct {
	cacheLock sync.Mutex
	cache     map[string]Release
}

func (r *ReleaseService) loadRelease(jsonPath string) error {
	releaseList := make([]Release, 0)
	// 读取json文件
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return common.NewError("读取 xray 版本文件失败: %v", err)
	}
	err = json.Unmarshal(data, &releaseList)
	if err != nil {
		return common.NewError("解析 xray 版本文件失败: %v", err)
	}

	for _, release := range releaseList {
		r.cacheLock.Lock()
		r.cache[release.TagName] = release
		r.cacheLock.Unlock()
	}
	return nil
}

func (r *ReleaseService) GetReleaseByTag(tagName string) (*Release, error) {
	if release, ok := r.cache[tagName]; ok {
		return &release, nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/XTLS/Xray-core/releases/tags/%v", tagName)
	resp, err := http.Get(url)
	if err != nil {
		return &Release{}, err
	}

	defer resp.Body.Close()
	buffer := bytes.NewBuffer(make([]byte, 8192))
	buffer.Reset()
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return &Release{}, err
	}

	release := Release{}
	err = json.Unmarshal(buffer.Bytes(), &release)
	if err != nil {
		return &Release{}, err
	}

	r.cacheLock.Lock()
	r.cache[tagName] = release
	r.cacheLock.Unlock()

	return &release, nil
}

type Release struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		UserViewType      string `json:"user_view_type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
	Reactions  struct {
		URL        string `json:"url"`
		TotalCount int    `json:"total_count"`
		Num1       int    `json:"+1"`
		Num10      int    `json:"-1"`
		Laugh      int    `json:"laugh"`
		Hooray     int    `json:"hooray"`
		Confused   int    `json:"confused"`
		Heart      int    `json:"heart"`
		Rocket     int    `json:"rocket"`
		Eyes       int    `json:"eyes"`
	} `json:"reactions"`
}
