import { useState, useEffect } from 'react';

export const usePlaylists = () => {
  const [playlists, setPlaylists] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchPlaylists = async () => {
    try {
      const response = await fetch('music/playlists/');
      if (!response.ok) throw new Error('Failed to fetch playlists');
      const data = await response.json();
      setPlaylists(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const createPlaylist = async (playlistData) => {
    const response = await fetch('music/playlists/', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(playlistData),
    });
    
    if (!response.ok) throw new Error('Failed to create playlist');
    const newPlaylist = await response.json();
    setPlaylists([...playlists, newPlaylist]);
    return newPlaylist;
  };

  useEffect(() => {
    fetchPlaylists();
  }, []);

  return { playlists, loading, error, createPlaylist, refreshPlaylists: fetchPlaylists };
};

export const useRecentlyPlayed = () => {
  const [recentTracks, setRecentTracks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchRecentlyPlayed = async () => {
    try {
      const response = await fetch('music/recently-played/');
      if (!response.ok) throw new Error('Failed to fetch recently played');
      const data = await response.json();
      setRecentTracks(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const addRecentlyPlayed = async (songId) => {
    const response = await fetch(`/music/recently-played/${songId}`, {
      method: 'POST',
    });
    
    if (!response.ok) throw new Error('Failed to add to recently played');
    await fetchRecentlyPlayed();
  };

  useEffect(() => {
    fetchRecentlyPlayed();
  }, []);

  return { recentTracks, loading, error, addRecentlyPlayed };
};

export const useSearch = () => {
  const [searchResults, setSearchResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const searchMusic = async (query) => {
    if (!query.trim()) {
      setSearchResults([]);
      return;
    }

    setLoading(true);
    try {
      const response = await fetch(`/music/search/?q=${encodeURIComponent(query)}`);
      if (!response.ok) throw new Error('Search failed');
      const data = await response.json();
      setSearchResults(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return { searchResults, loading, error, searchMusic };
};