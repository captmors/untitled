import { Play } from 'lucide-react';

const RecentlyPlayed = () => {
  const recentTracks = [
    { id: 1, name: "Song Name 1", artist: "Artist 1", duration: "3:45", imageUrl: "/api/placeholder/50/50" },
    { id: 2, name: "Song Name 2", artist: "Artist 2", duration: "4:20", imageUrl: "/api/placeholder/50/50" },
    { id: 3, name: "Song Name 3", artist: "Artist 3", duration: "3:15", imageUrl: "/api/placeholder/50/50" },
    { id: 4, name: "Song Name 4", artist: "Artist 4", duration: "2:55", imageUrl: "/api/placeholder/50/50" },
    { id: 5, name: "Song Name 5", artist: "Artist 5", duration: "3:30", imageUrl: "/api/placeholder/50/50" },
  ];

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
                src={track.imageUrl} 
                alt={track.name}
                className="w-10 h-10 rounded mr-4"
              />
              <div>
                <div className="text-white">{track.name}</div>
                <div className="text-gray-400 text-sm">{track.artist}</div>
              </div>
            </div>
            <div className="w-32 flex items-center justify-end space-x-4">
              <button className="invisible group-hover:visible">
                <Play className="h-4 w-4 text-white" />
              </button>
              <span className="text-gray-400">{track.duration}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default RecentlyPlayed;