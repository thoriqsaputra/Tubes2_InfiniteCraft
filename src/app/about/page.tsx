import Button from '../ui/components/Button';
import Link from 'next/link';
import Image from 'next/image';

export default function Page() {
  return (
    <div className="bg-gradient-to-r from-fuchsia-500 to-sky-500 min-h-screen">
    <div className="flex justify-normal px-1 pt-2 md:px-16 md:py-4">
      <div className="flex space-x-4">
        <Button href="/">Home</Button>
        <Button href="/how-to-use">How to Use</Button>
      </div>
    </div>

    <div className="max-w-4xl mx-auto px-8 py-10 flex flex-col md:flex-row justify-center items-center">
      <div className="md:w-1/2 text-center md:text-left">
        <h1 className="text-4xl font-bold text-white mb-4">Explore Wikipedia Like Never Before</h1>
        <p className="text-xl text-white opacity-75 leading-relaxed mb-8">
          We leverage the power of pathfinding algorithms to unlock the vast network of knowledge on Wikipedia. Dive deep into any topic with ease, and uncover hidden connections you never knew existed.
        </p>
        <Link href="/" className="bg-indigo-500 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded-full shadow-md transition duration-300">Start Exploring</Link>
      </div>
      <div className="md:w-1/2 flex justify-center items-center">
        <img src="/wikipedia.png" alt="Wikipedia" className="w-full max-w-md" />
      </div>
    </div>

    <div className="grid grid-cols-1 md:grid-cols-3 gap-8 px-16 py-8">
      {/* Creator Card 1 */}
      <div className="bg-white rounded-lg shadow-md p-6 transform hover:scale-105 transition-transform duration-300">
        <h3 className="text-xl font-semibold text-orange-500 mb-2">Christian Justin Hendrawan</h3>
        <p className="text-gray-700">Create IDS Algorithm for a deep path.</p>
      </div>
      {/* Creator Card 2 */}
      <div className="bg-white rounded-lg shadow-md p-6 transform hover:scale-105 transition-transform duration-300">
        <h3 className="text-xl font-semibold text-emerald-500 mb-2">Ahmad Thoriq Saputra</h3>
        <p className="text-gray-700">leads our development efforts, specializing in web development.</p>
      </div>
      {/* Creator Card 3 */}
      <div className="bg-white rounded-lg shadow-md p-6 transform hover:scale-105 transition-transform duration-300">
        <h3 className="text-xl font-semibold text-violet-500 mb-2">M. Zaidan Sa'dun </h3>
        <p className="text-gray-700">Implementing an effiecient BFS Algorithm to find the shortest path.</p>
      </div>
    </div>
  </div>
  );
}