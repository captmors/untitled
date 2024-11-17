import { Search, Bell, User } from 'lucide-react';
import { useState } from 'react';
import { useAuth } from '../hooks/useAuth';
import { useSearch } from '../hooks/useMusic';

const Header = () => {
  const { user, logout } = useAuth();
  const { searchResults, loading, searchMusic } = useSearch();
  const [showDropdown, setShowDropdown] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearch = (e) => {
    const query = e.target.value;
    setSearchQuery(query);
    searchMusic(query);
  };

  return (
    <header className="fixed top-0 left-0 right-0 bg-gray-900 text-white z-50">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
        <div className="flex items-center space-x-4">
          <a href="/" className="text-xl font-bold">
            MusicApp
          </a>
        </div>

          <div className="flex-1 max-w-xl mx-4">
            <div className="relative">
              <input
                type="text"
                value={searchQuery}
                onChange={handleSearch}
                placeholder="Search for songs, artists, or playlists..."
                className="w-full bg-gray-800 rounded-full py-2 px-4 pl-10 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <Search className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
              
              {/* Search Results Dropdown */}
              {searchQuery && (
                <div className="absolute w-full bg-gray-800 mt-2 rounded-lg shadow-lg">
                  {loading ? (
                    <div className="p-4 text-gray-400">Searching...</div>
                  ) : (
                    searchResults.map((result) => (
                      <div
                        key={result.id}
                        className="p-2 hover:bg-gray-700 cursor-pointer"
                      >
                        {result.name}
                      </div>
                    ))
                  )}
                </div>
              )}
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <button className="p-2 hover:bg-gray-800 rounded-full">
              <Bell className="h-5 w-5" />
            </button>
            <div className="relative">
              <button 
                className="p-2 hover:bg-gray-800 rounded-full"
                onClick={() => setShowDropdown(!showDropdown)}
              >
                <User className="h-5 w-5" />
              </button>
              
              {showDropdown && (
                <div className="absolute right-0 mt-2 w-48 bg-gray-800 rounded-lg shadow-lg">
                  {user ? (
                    <>
                      <div className="px-4 py-2 border-b border-gray-700">
                        <div className="font-medium">{user.name}</div>
                        <div className="text-sm text-gray-400">{user.email}</div>
                      </div>
                      <button
                        onClick={logout}
                        className="w-full text-left px-4 py-2 hover:bg-gray-700"
                      >
                        Logout
                      </button>
                    </>
                  ) : (
                    <div className="px-4 py-2">
                      <a href="/login" className="block py-1 hover:text-blue-400">Login</a>
                      <a href="/register" className="block py-1 hover:text-blue-400">Register</a>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;