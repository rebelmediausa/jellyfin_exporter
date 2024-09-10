// Copyright 2010 Rebel Media
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nomedia
// +build !nomedia

package collector

import (
	"errors"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rebelmediausa/jellyfin_exporter/collector/utils"
	"github.com/rebelmediausa/jellyfin_exporter/config"
)

var (
	mediaMoviesValue      float64
	mediaShowsValue       float64
	mediaEpisodesValue    float64
	mediaArtistsValue     float64
	mediaProgramsValue    float64
	mediaTrailersValue    float64
	mediaSongsValue       float64
	mediaAlbumsValue      float64
	mediaMusicVideosValue float64
	mediaBoxSetsValue     float64
	mediaBooksValue       float64
	mediaItemsValue       float64
)

type mediaCollector struct {
	mediaMovies      *prometheus.Desc
	mediaShows       *prometheus.Desc
	mediaEpisodes    *prometheus.Desc
	mediaArtists     *prometheus.Desc
	mediaPrograms    *prometheus.Desc
	mediaTrailers    *prometheus.Desc
	mediaSongs       *prometheus.Desc
	mediaAlbums      *prometheus.Desc
	mediaMusicVideos *prometheus.Desc
	mediaBoxSets     *prometheus.Desc
	mediaBooks       *prometheus.Desc
	mediaItems       *prometheus.Desc
	logger           log.Logger
}

func init() {
	registerCollector("media", defaultEnabled, NewMediaCollector)
}

func NewMediaCollector(logger log.Logger) (Collector, error) {
	const subsystem = "media"
	mediaMovies := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "movies_count",
		),
		"Total Movies.",
		[]string{}, nil,
	)
	mediaShows := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "shows_count",
		),
		"Total TV Shows.",
		[]string{}, nil,
	)
	mediaEpisodes := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "episodes_count",
		),
		"Total TV Show Episodes.",
		[]string{}, nil,
	)
	mediaArtists := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "artists_count",
		),
		"Total Artists.",
		[]string{}, nil,
	)
	mediaPrograms := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "programs_count",
		),
		"Total Programs.",
		[]string{}, nil,
	)
	mediaTrailers := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "trailers_count",
		),
		"Total Trailers.",
		[]string{}, nil,
	)
	mediaSongs := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "songs_count",
		),
		"Total Songs.",
		[]string{}, nil,
	)
	mediaAlbums := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "albums_count",
		),
		"Total Albums.",
		[]string{}, nil,
	)
	mediaMusicVideos := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "music_videos_count",
		),
		"Total Music Videos.",
		[]string{}, nil,
	)
	mediaBoxSets := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "box_sets_count",
		),
		"Total Box Sets.",
		[]string{}, nil,
	)
	mediaBooks := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "books_count",
		),
		"Total Books.",
		[]string{}, nil,
	)
	mediaItems := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, subsystem, "items_count",
		),
		"Total Items.",
		[]string{}, nil,
	)
	return &mediaCollector{
		mediaMovies:      mediaMovies,
		mediaShows:       mediaShows,
		mediaEpisodes:    mediaEpisodes,
		mediaArtists:     mediaArtists,
		mediaPrograms:    mediaPrograms,
		mediaTrailers:    mediaTrailers,
		mediaSongs:       mediaSongs,
		mediaAlbums:      mediaAlbums,
		mediaMusicVideos: mediaMusicVideos,
		mediaBoxSets:     mediaBoxSets,
		mediaBooks:       mediaBooks,
		mediaItems:       mediaItems,
		logger:           logger,
	}, nil
}

func (c *mediaCollector) Update(ch chan<- prometheus.Metric) error {
	jellyfinURL, jellyfinToken, nil := config.JellyfinInfo(c.logger)

	jellyfinAPIURL := fmt.Sprintf("%s/Items/Counts", jellyfinURL)
	rawData, statusCode, err := utils.GetHTTP(jellyfinAPIURL, jellyfinToken)
	if !errors.Is(err, nil) {
		level.Error(c.logger).Log("msg", "Error fetching API:", "err", err)
		return err
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		level.Error(c.logger).Log("msg", "Error parsing Media response", "err", err)
		return err
	}

	if statusCode == 200 {
		mediaMoviesValue = data["MovieCount"].(float64)
		mediaShowsValue = data["SeriesCount"].(float64)
		mediaEpisodesValue = data["EpisodeCount"].(float64)
		mediaArtistsValue = data["ArtistCount"].(float64)
		mediaArtistsValue = data["ArtistCount"].(float64)
		mediaProgramsValue = data["ProgramCount"].(float64)
		mediaTrailersValue = data["TrailerCount"].(float64)
		mediaSongsValue = data["SongCount"].(float64)
		mediaAlbumsValue = data["AlbumCount"].(float64)
		mediaMusicVideosValue = data["MusicVideoCount"].(float64)
		mediaBoxSetsValue = data["BoxSetCount"].(float64)
		mediaBooksValue = data["BookCount"].(float64)
		mediaItemsValue = data["ItemCount"].(float64)
	}

	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Movies", mediaMoviesValue)
	ch <- prometheus.MustNewConstMetric(c.mediaMovies, prometheus.CounterValue, mediaMoviesValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "TV_Shows", mediaShowsValue)
	ch <- prometheus.MustNewConstMetric(c.mediaShows, prometheus.CounterValue, mediaShowsValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "TV_Show_Episodes", mediaEpisodesValue)
	ch <- prometheus.MustNewConstMetric(c.mediaEpisodes, prometheus.CounterValue, mediaEpisodesValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Artists", mediaArtistsValue)
	ch <- prometheus.MustNewConstMetric(c.mediaArtists, prometheus.CounterValue, mediaArtistsValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Programs", mediaProgramsValue)
	ch <- prometheus.MustNewConstMetric(c.mediaPrograms, prometheus.CounterValue, mediaProgramsValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Trailers", mediaTrailersValue)
	ch <- prometheus.MustNewConstMetric(c.mediaTrailers, prometheus.CounterValue, mediaTrailersValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Songs", mediaSongsValue)
	ch <- prometheus.MustNewConstMetric(c.mediaSongs, prometheus.CounterValue, mediaSongsValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Albums", mediaAlbumsValue)
	ch <- prometheus.MustNewConstMetric(c.mediaAlbums, prometheus.CounterValue, mediaAlbumsValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Music_Videos", mediaMusicVideosValue)
	ch <- prometheus.MustNewConstMetric(c.mediaMusicVideos, prometheus.CounterValue, mediaMusicVideosValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Box_Sets", mediaBoxSetsValue)
	ch <- prometheus.MustNewConstMetric(c.mediaBoxSets, prometheus.CounterValue, mediaBoxSetsValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Books", mediaBooksValue)
	ch <- prometheus.MustNewConstMetric(c.mediaBooks, prometheus.CounterValue, mediaBooksValue)
	level.Debug(c.logger).Log("msg", "Jellyfin Media System Total", "Items", mediaItemsValue)
	ch <- prometheus.MustNewConstMetric(c.mediaItems, prometheus.CounterValue, mediaItemsValue)

	return nil
}
