import Header from '../components/Header';
import PlaylistCarousel from '../components/PlaylistCarousel';
import RecentlyPlayed from '../components/RecentlyPlayed';

function Home() {
  return (
    <div className="home-page">
      <Header />
      <PlaylistCarousel title="Featured Playlists" />
      <RecentlyPlayed />
    </div>
  );
}

export default Home;
