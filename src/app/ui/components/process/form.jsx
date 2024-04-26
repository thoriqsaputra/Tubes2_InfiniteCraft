"use client";

import  Input  from './Input';
import Image from 'next/image';
import RadioMethod from './RadioMethod';
import RadioLanguage from './RadioLanguage';
import { useState, useEffect } from 'react';
import RadioSolution from './RadioSolution';


export default function ProcessForm() {
    const [startValue, setStartValue] = useState("");
    const [destinationValue, setDestinationValue] = useState("");
    const [disabled, setDisabled] = useState(true);
    const [clicked, setClicked] = useState(false);
    const [selectedMethod, setSelectedMethod] = useState(null);
    const [selectedLanguage, setSelectedLanguage] = useState(null);
    const [disableLanguage, setDisableLanguage] = useState(true);
    const [selectedSolution, setSelectedSolution] = useState(null);

    useEffect(() => {
        setDisabled(startValue.trim() === "" || selectedMethod === null || selectedLanguage === null);
    }, [selectedMethod, selectedLanguage]);

    const handleSolutionChange = (solution) => {
        setSelectedSolution(solution);
    };

    const handleMethodChange = (method) => {
        setSelectedMethod(method);
    };

    const handleLanguageChange = (language) => {
        setSelectedLanguage(language);
        setDisableLanguage(false);
    };

    const handleChangeDestination = (e) => {
        const newValue = e.target.value;
        setDestinationValue(newValue);
        setDisabled(newValue.trim() === "" || startValue.trim() === "" || selectedMethod === null || selectedLanguage === null);
    };
    
    const handleChangeStart = (e) => {
        const newValue = e.target.value;
        setStartValue(newValue);
        setDisabled(newValue.trim() === "" || destinationValue.trim() === "" || selectedMethod === null || selectedLanguage === null);
    };  

    const handleClick = (e) => {
        setClicked(true);
        setTimeout(() => {
          setClicked(false);
        }, 300);
        setDisabled(true);
        e.preventDefault();
        const audio = new Audio("/vroom.mp3");
        audio.play();
        audio.onended = () => {
            setDisabled(false);
        };
    };

  return (
    <div className='w-screen pt-12'>
        <div className="flex flex-col md:flex-col md:gap-11 lg:flex-row justify-center items-center">
           <div><Input placeholder="start" handleChange={handleChangeStart} inputValue={startValue} language={selectedLanguage} disable={disableLanguage}/></div>
            <div className='py-10 md:py-10 lg:py-0'>
                <Image
                    src="/arrowRight.png"
                    alt="arrow"
                    width={100}
                    height={100}
                    className='transform rotate-90 h-auto w-auto md:rotate-90 lg:rotate-0'
                />
            </div>
            <div><Input placeholder="destination" handleChange={handleChangeDestination} inputValue={destinationValue} language={selectedLanguage} disable={disableLanguage}/></div>
        </div>
        <div className='flex justify-center lg:gap-5 md:gap-0 items-center lg:pt-[80px] pt-12'>
            <div className='md:mt-[-30px] mt-0'>
                <RadioLanguage handleOptionChange={handleLanguageChange} selectedOption={selectedLanguage} />
            </div>
            <div>
                <button
                disabled={disabled}
                className={`bg-gradient-to-br from-amber-500 to-red-400 outline-2 outline shadow-lg text-white text-xl md:text-3xl py-1 px-2 font-semibold md:py-2 md:px-4 rounded-lg transition ease-in-out duration-500 
                ${clicked ? 'scale-95 ' : ''} hover:from-green-400 hover:to-blue-400 focus:outline-none`}
                onClick={handleClick}
                >
                    Start Racing!!
                </button>
            </div>
            <div className='md:mt-[-30px] md:mx-6 mt-0'>
                <RadioMethod handleOptionChange={handleMethodChange} selectedOption={selectedMethod}/>
            </div>
        </div>
        <RadioSolution handleOptionChange={handleSolutionChange} selectedOption={selectedSolution}/>
    </div>
  );
}
