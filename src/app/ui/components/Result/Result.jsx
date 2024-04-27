'use client';

import Image from 'next/image';
import Link from 'next/link';
import { useEffect } from 'react';

export default function Result({Result, showLoading, method}) {
    const pedro = new Audio('pedro.mp3');
    const smooth = new Audio('smooth.mp3');
    const article = Result.path;

    // const [articleDetails, setArticleDetails] = useState([]);
    
    // function constructWikipediaURL(language, inputValue) {
    //     const wikipediaEndpoint = `https://${language}.wikipedia.org/w/api.php`;
    //     const wikipediaParams =
    //         `?action=query` +
    //         `&format=json` +
    //         `&gpssearch=${inputValue}` +
    //         `&generator=prefixsearch` +
    //         `&prop=pageprops%7Cpageimages%7Cpageterms` +
    //         `&redirects=` +
    //         `&ppprop=displaytitle` +
    //         `&piprop=thumbnail` +
    //         `&pithumbsize=160` +
    //         `&pilimit=30` +
    //         `&wbptterms=description` +
    //         `&gpsnamespace=0` +
    //         `&gpslimit=5` +
    //         `&origin=*`;
    
    //     return wikipediaEndpoint + wikipediaParams;
    // }

    // useEffect(() => {
    //     const fethDetails = async () => {
    //         const detailPromises = Result.path.map(async (path) => {

    // }, [Result]);
    

    useEffect(() => {
        if (showLoading) {
            if(method === "BFS") {
                pedro.loop = true;
                pedro.play();
            } else {
                smooth.loop = true;
                smooth.play();
            }
        } else {
            if(method === "IDS") {
                pedro.pause();
                pedro.currentTime = 0;
            } else {
                smooth.pause();
                smooth.currentTime = 0;
            }
        }
    }, [showLoading]);

    return (
        <main className='min-h-screen bg-gradient-radial from-amber-600 to-amber-300 relative'>
            <div className="curve">
                <svg data-name="Layer 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 120" preserveAspectRatio="none">
                <path d="M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z" opacity=".25" class="shape-fill"></path>
                <path d="M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z" opacity=".5" class="shape-fill"></path>
                <path d="M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z" class="shape-fill"></path>
            </svg>
            </div>
            <div className="min-h-screen flex items-center justify-center">
                {showLoading ? (
                    <Image alt="loading" src={method === "IDS" ? '/f1.gif' : '/mclaren.gif'} width={400} height={400}
                    className='rounded-full' />
                ) : 
                (
                <div className="flex flex-col gap-8 justify-center items-center px-4 py-[100px]">
                <div className="bg-blue-600 text-white rounded-lg shadow-md p-8 flex-col items-center transition-all duration-300 hover:shadow-lg">
                    <h2 className="text-2xl font-bold mb-4">Path has been found!</h2>
                    <p>Time Taken: {Result.duration} ms</p>
                    <p>Articles Checked: {Result.links}</p>
                    <p>Degrees: {Result.degrees}</p>
                </div>
                <div className="grid grid-cols-1 gap-10 items-center justify-center md:grid-cols-2 lg:grid-cols-3 mt-8">
                    {article.map((list, index) => (
                        <div key={index} className="bg-white min-w-[300px] outline outline-offset-2 outline-2 outline-red-600 rounded-lg shadow-md flex-col justify-center transition-all duration-300 hover:shadow-lg hover:shadow-sky-600">
                            <h2 className="text-2xl text-center text-black font-bold my-4">Path {index + 1}</h2>
                            <ul className=''>
                                {list.map((article, idx) => (
                                    <li key={idx} className="text-blue-600 hover:text-white py-4 px-3 rounded-xl hover:bg-orange-400 transition-all duration-300 transform hover:scale-105">
                                    <Link href={`https://${Result.language}.wikipedia.org/wiki/${encodeURIComponent(article)}`} className="block w-full h-full">
                                        <div className='flex flex-row justify-normal items-center gap-10'>
                                            {/* <Image src={article.image} alt="thumbnail" width={50} height={50} className="rounded-md" /> */}
                                            <span>{article}</span>
                                        </div>
                                    </Link>
                                </li>
                                ))}
                            </ul>
                        </div>
                    ))}
                </div>
                </div>)}
            </div>
            <div class="wave">
                <svg data-name="Layer 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 120" preserveAspectRatio="none">
                    <path d="M321.39,56.44c58-10.79,114.16-30.13,172-41.86,82.39-16.72,168.19-17.73,250.45-.39C823.78,31,906.67,72,985.66,92.83c70.05,18.48,146.53,26.09,214.34,3V0H0V27.35A600.21,600.21,0,0,0,321.39,56.44Z" class="shape-fill"></path>
                </svg>
            </div>
        </main>
    )
}

