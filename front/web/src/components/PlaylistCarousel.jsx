import { ChevronLeft, ChevronRight, Play } from 'lucide-react';
import PropTypes from 'prop-types';

const PlaylistCarousel = ({ title }) => {
  const playlists = [
    { id: 1, name: "Today's Hits", imageUrl: "/api/placeholder/200/200", songs: "50 songs" },
    { id: 2, name: "Chill Vibes", imageUrl: "/api/placeholder/200/200", songs: "45 songs" },
    { id: 3, name: "Rock Classics", imageUrl: "/api/placeholder/200/200", songs: "75 songs" },
    { id: 4, name: "Jazz Essentials", imageUrl: "/api/placeholder/200/200", songs: "30 songs" },
    { id: 5, name: "Workout Mix", imageUrl: "/api/placeholder/200/200", songs: "60 songs" },
  ];

  return (
    <div className="mt-20 px-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-2xl font-bold text-white">{title}</h2>
        <div className="flex space-x-2">
          <button className="p-2 rounded-full bg-gray-800 hover:bg-gray-700">
            <ChevronLeft className="h-5 w-5 text-white" />
          </button>
          <button className="p-2 rounded-full bg-gray-800 hover:bg-gray-700">
            <ChevronRight className="h-5 w-5 text-white" />
          </button>
        </div>
      </div>
      
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
        {playlists.map((playlist) => (
          <div 
            key={playlist.id} 
            className="bg-gray-800 rounded-lg p-4 hover:bg-gray-700 transition-colors"
          >
            <div className="relative group">
              <img 
                src={playlist.imageUrl} 
                alt={playlist.name}
                className="w-full aspect-square object-cover rounded-md"
              />
              <button className="absolute bottom-2 right-2 p-3 bg-green-500 rounded-full opacity-0 group-hover:opacity-100 transition-opacity">
                <Play className="h-6 w-6 text-white" fill="white" />
              </button>
            </div>
            <h3 className="mt-2 text-white font-semibold">{playlist.name}</h3>
            <p className="text-gray-400 text-sm">{playlist.songs}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

PlaylistCarousel.propTypes = {
    title: PropTypes.string.isRequired,
  };

export default PlaylistCarousel;