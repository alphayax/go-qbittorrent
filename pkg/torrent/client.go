package torrent

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/sirupsen/logrus"

	"github.com/KnutZuidema/go-qbittorrent/pkg"
	"github.com/KnutZuidema/go-qbittorrent/pkg/model"
)

type Client struct {
	BaseUrl string
	Client  *http.Client
	Logger  logrus.FieldLogger
}

func (c Client) GetList(options *model.GetTorrentListOptions) ([]*model.Torrent, error) {
	endpoint := c.BaseUrl + "/info"
	if options != nil {
		params, err := query.Values(options)
		if err != nil {
			return nil, err
		}
		endpoint += "?" + params.Encode()
	}
	var res []*model.Torrent
	if err := pkg.GetInto(c.Client, &res, endpoint, nil); err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) GetProperties(hash string) (*model.TorrentProperties, error) {
	params := url.Values{}
	params.Add("hash", hash)
	endpoint := c.BaseUrl + "/properties?" + params.Encode()
	var res model.TorrentProperties
	if err := pkg.GetInto(c.Client, &res, endpoint, nil); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c Client) GetTrackers(hash string) ([]*model.TorrentTracker, error) {
	params := url.Values{}
	params.Add("hash", hash)
	endpoint := c.BaseUrl + "/trackers?" + params.Encode()
	var res []*model.TorrentTracker
	if err := pkg.GetInto(c.Client, &res, endpoint, nil); err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) GetWebSeeds(hash string) ([]string, error) {
	params := url.Values{}
	params.Add("hash", hash)
	endpoint := c.BaseUrl + "/trackers?" + params.Encode()
	var res []struct {
		URL string `json:"url"`
	}
	if err := pkg.GetInto(c.Client, &res, endpoint, nil); err != nil {
		return nil, err
	}
	var seeds []string
	for _, seed := range res {
		seeds = append(seeds, seed.URL)
	}
	return seeds, nil
}

func (c Client) GetContents(hash string) ([]*model.TorrentContent, error) {
	params := url.Values{}
	params.Add("hash", hash)
	endpoint := c.BaseUrl + "/files?" + params.Encode()
	var res []*model.TorrentContent
	if err := pkg.GetInto(c.Client, &res, endpoint, nil); err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) GetPieceStates(hash string) ([]model.TorrentPieceState, error) {
	params := url.Values{}
	params.Add("hash", hash)
	endpoint := c.BaseUrl + "/pieceStates?" + params.Encode()
	var res []model.TorrentPieceState
	if err := pkg.GetInto(c.Client, &res, endpoint, nil); err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) GetPieceHashes(hash string) ([]string, error) {
	params := url.Values{}
	params.Add("hash", hash)
	endpoint := c.BaseUrl + "/pieceHashes?" + params.Encode()
	var res []string
	if err := pkg.GetInto(c.Client, &res, endpoint, nil); err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) StopTorrents(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/pause", strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) ResumeTorrents(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/resume", strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) DeleteTorrents(hashes []string, deleteFiles bool) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	params.Add("deleteFiles", fmt.Sprintf("%v", deleteFiles))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/delete", strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) RecheckTorrents(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	endpoint := c.BaseUrl + "/recheck?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) ReannounceTorrents(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	endpoint := c.BaseUrl + "/reannounce?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) AddURLs(urls []string, options *model.AddTorrentsOptions) error {
	if err := pkg.PostMultipartLinks(c.Client, c.BaseUrl+"/add", options, urls); err != nil {
		return err
	}
	return nil
}

func (c Client) AddFiles(files map[string][]byte, options *model.AddTorrentsOptions) error {
	if err := pkg.PostMultipartFiles(c.Client, c.BaseUrl+"/add", options, files); err != nil {
		return err
	}
	return nil
}

func (c Client) AddTrackers(hash string, trackerURLs []string) error {
	params := url.Values{}
	params.Add("hash", hash)
	params.Add("urls", strings.Join(trackerURLs, "\n"))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/addTrackers",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) EditTrackers(hash, old, new string) error {
	params := url.Values{}
	params.Add("hash", hash)
	params.Add("origUrl", old)
	params.Add("newUrl", new)
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/editTracker",
		strings.NewReader(params.Encode()),
		"application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) RemoveTrackers(hash string, trackerURLs []string) error {
	params := url.Values{}
	params.Add("hash", hash)
	params.Add("urls", strings.Join(trackerURLs, "|"))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/removeTrackers",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) IncreasePriority(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	endpoint := c.BaseUrl + "/increasePrio?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) DecreasePriority(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	endpoint := c.BaseUrl + "/decreasePrio?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) SetMaximumPriority(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	endpoint := c.BaseUrl + "/topPrio?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) SetMinimumPriority(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	endpoint := c.BaseUrl + "/bottomPrio?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) SetFilePriorities(hash string, ids []string, priority model.TorrentPriority) error {
	params := url.Values{}
	params.Add("hash", hash)
	params.Add("id", strings.Join(ids, "|"))
	params.Add("priority", strconv.Itoa(int(priority)))
	endpoint := c.BaseUrl + "/filePrio?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) GetDownloadLimits(hashes []string) (map[string]int, error) {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	var res map[string]int
	if err := pkg.GetIntoWithContentType(c.Client, &res, c.BaseUrl+"/downloadLimit",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) SetDownloadLimits(hashes []string, limit int) error {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	params.Add("limit", strconv.Itoa(limit))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/setDownloadLimit",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) SetShareLimits(hashes []string, ratioLimit float64, seedingTimeLimit int) error {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	params.Add("ratioLimit", strconv.FormatFloat(ratioLimit, 'f', -1, 64))
	params.Add("seedingTimeLimit", strconv.Itoa(seedingTimeLimit))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/setShareLimits",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) GetUploadLimits(hashes []string) (map[string]int, error) {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	var res map[string]int
	if err := pkg.GetIntoWithContentType(c.Client, &res, c.BaseUrl+"/uploadLimit",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) SetUploadLimits(hashes []string, limit int) error {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	params.Add("limit", strconv.Itoa(limit))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/setUploadLimit",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) SetLocations(hashes []string, location string) error {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	params.Add("location", location)
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/setLocation",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) SetName(hash string, name string) error {
	params := url.Values{}
	params.Add("hash", hash)
	params.Add("name", name)
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/rename",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) SetCategories(hashes []string, category string) error {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	params.Add("category", category)
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/setCategory",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) GetCategories() (map[string]*model.Category, error) {
	var res map[string]*model.Category
	if err := pkg.GetInto(c.Client, &res, c.BaseUrl+"/categories", nil); err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) AddCategory(category string, savePath string) error {
	params := url.Values{}
	params.Add("category", category)
	params.Add("savePath", savePath)
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/createCategory",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) EditCategory(category string, savePath string) error {
	params := url.Values{}
	params.Add("category", category)
	params.Add("savePath", savePath)
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/editCategory",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) RemoveCategory(categories []string) error {
	params := url.Values{}
	params.Add("categories", strings.Join(categories, "\n"))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/removeCategories",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) SetAutomaticManagement(hashes []string, enable bool) error {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	params.Add("enable", fmt.Sprintf("%v", enable))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/setAutoManagement",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) ToggleSequentialDownload(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	endpoint := c.BaseUrl + "/toggleSequentialDownload?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) ToggleFirstLastPiecePriority(hashes []string) error {
	value := strings.Join(hashes, "|")
	params := url.Values{}
	params.Add("hashes", value)
	endpoint := c.BaseUrl + "/toggleFirstLastPiecePrio?" + params.Encode()
	if err := pkg.Post(c.Client, endpoint, nil); err != nil {
		return err
	}
	return nil
}

func (c Client) SetForceStart(hashes []string, enable bool) error {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	params.Add("value", fmt.Sprintf("%v", enable))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/setForceStart",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}

func (c Client) SetSuperSeeding(hashes []string, enable bool) error {
	params := url.Values{}
	params.Add("hashes", strings.Join(hashes, "|"))
	params.Add("value", fmt.Sprintf("%v", enable))
	if err := pkg.PostWithContentType(c.Client, c.BaseUrl+"/setSuperSeeding",
		strings.NewReader(params.Encode()), "application/x-www-form-urlencoded"); err != nil {
		return err
	}
	return nil
}
