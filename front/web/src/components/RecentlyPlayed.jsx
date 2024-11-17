import { Play } from 'lucide-react';
import { useRecentlyPlayed } from '../hooks/useMusic';

const RecentlyPlayed = () => {
  const { recentTracks, loading, error, addRecentlyPlayed } = useRecentlyPlayed();

  if (loading) {
    return (
      <div className="px-4 mt-8">
        <h2 className="text-2xl font-bold text-white">Recently Played</h2>
        <div className="mt-4 text-gray-400">Loading recently played tracks...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="px-4 mt-8">
        <h2 className="text-2xl font-bold text-white">Recently Played</h2>
        <div className="mt-4 text-red-500">Error loading tracks: {error}</div>
      </div>
    );
  }

  const handlePlay = async (trackId) => {
    try {
      await addRecentlyPlayed(trackId);
      // TODO play logic 
    } catch (err) {
      console.error('Failed to update recently played:', err);
    }
  };

  return (
    <div className="px-4 mt-8">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-2xl font-bold text-white">Recently Played</h2>
        <button className="text-gray-400 hover:text-white">See All</button>
      </div>

      <div className="bg-gray-800 rounded-lg">
        {recentTracks.map((track, index) => (
          <div 
            key={track.id}
            className="flex items-center px-4 py-2 hover:bg-gray-700 group"
          >
            <div className="w-8 text-gray-400">{index + 1}</div>
            <div className="flex-1 flex items-center">
              <img 
                src={track.imageUrl || "/api/placeholder/50/50"} 
                alt={track.name}
                className="w-10 h-10 rounded mr-4"
              />
              <div>
                <div className="text-white">{track.name}</div>
                <div className="text-gray-400">{track.artist}</div>
              </div>
            </div>
            <button 
              onClick={() => handlePlay(track.id)}
              className="ml-4 text-gray-400 hover:text-white group-hover:text-white"
            >
              <Play size={20} />
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default RecentlyPlayed;
