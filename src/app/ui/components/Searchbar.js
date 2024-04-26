"use client";

import { useState } from 'react';

const SearchBar = () => {
    const [query, setQuery] = useState('');
    const [searchResults, setSearchResults] = useState([]);

    const handleInputChange = async (event) => {
        const newQuery = event.target.value;
        setQuery(newQuery);

        if (newQuery.trim() === '') {
            setSearchResults([]);
            return;
        }

        try {
            const response = await fetch(`https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=${encodeURIComponent(newQuery)}`);
            const data = await response.json();
            const results = data.query.search.map(result => ({ title: result.title }));
            setSearchResults(results);
        } catch (error) {
            console.error('Error fetching data from Wikipedia:', error);
            setSearchResults([]);
        }
    };

    const handleResultClick = (title) => {
        setQuery(title);
        setSearchResults([]);
    };

    return (
        <div className="search-container relative">
            <input
                type="text"
                value={query}
                onChange={handleInputChange}
                placeholder="Search Wikipedia"
                className="w-full px-4 py-2 text-black rounded border border-gray-300 focus:outline-none focus:border-blue-400"
            />
            {searchResults.length > 0 && (
                <ul className="search-results absolute left-0 w-full bg-white border border-gray-300 rounded mt-1 shadow-md">
                    {searchResults.map((result, index) => (
                        <li
                            key={index}
                            onClick={() => handleResultClick(result.title)}
                            className="py-2 px-4 cursor-pointer hover:bg-gray-100"
                        >
                            {result.title}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default SearchBar;
