
"use client";

import React, { useEffect, useState } from "react";
import Image from "next/image";

const Input = ({ placeholder, handleChange, inputValue, language, disable }) => {
    const [searchResults, setSearchResults] = useState([]);
    const [showRecommendation, setShowRecommendation] = useState(false);
    const [recommendationClicked, setRecommendationClicked] = useState(false);

    const wikipediaEndpoint = 'https://' + language + '.wikipedia.org/w/api.php';
    const wikipediaParams =
        '?action=query' +
        '&format=json' +
        '&gpssearch=' + inputValue +
        '&generator=prefixsearch' +
        '&prop=pageprops%7Cpageimages%7Cpageterms' +
        '&redirects=' +
        '&ppprop=displaytitle' +
        '&piprop=thumbnail' +
        '&pithumbsize=160' +
        '&pilimit=30' +
        '&wbptterms=description' +
        '&gpsnamespace=0' +
        '&gpslimit=5' +
        '&origin=*';

    useEffect(() => {
        if (inputValue !== "" && !recommendationClicked) {
            fetch(wikipediaEndpoint + wikipediaParams)
                .then(response => response.json())
                .then(data => setSearchResults(data.query.pages))
                .catch(error => console.error('Error fetching data:', error));
            setShowRecommendation(true);
        } else {
            setSearchResults([]);
            setShowRecommendation(false);
        }
        setRecommendationClicked(false);
    }, [inputValue, language]);

    const handleRecommendationClick = (title) => {
        handleChange({ target: { value: title } });
        setShowRecommendation(false);
        setRecommendationClicked(true);
    };

    return (
        <main>
            <div>
                <input
                    type="text"
                    placeholder={placeholder}
                    value={inputValue}
                    onChange={handleChange}
                    disabled={disable}
                    className="text-black flex border-amber-[#ADBC9F] border-double border-4 rounded-md w-full min-w-[500px] px-3 py-2 focus:outline-1 focus:outline-[#627254] focus:border-[#76885B] hover:bg-[#FDF7E4] hover:border-[#FDEED1] hover:shadow-lg transition duration-500 ease-in-out md:min-w-[500px]"
                />
            </div>
            <div>
                {showRecommendation && (
                    <ul className="absolute w-[500px] bg-white border  border-gray-300 rounded mt-1 shadow-md max-h-[250px] overflow-y-auto">
                        {Object.keys(searchResults).map((key, index) => (
                            <li
                                key={index}
                                onClick={() => handleRecommendationClick(searchResults[key].title)}
                                className="py-3 px-4 cursor-pointer flex flex-row items-center text-black hover:bg-blue-200"
                            >
                                <div>
                                    <img
                                        src={searchResults[key].thumbnail ? searchResults[key].thumbnail.source : "/no-image.gif"}
                                        alt="thumbnail"
                                        width={60}
                                        height={60}
                                        className="rounded-md mr-2 "
                                    />
                                </div>
                                <div>
                                    {searchResults[key].title}
                                </div>
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </main>
    );
};

export default Input;