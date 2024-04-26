'use client';

import Image from "next/image";
import { racing, kotta } from "./ui/fonts";
import Button from "./ui/components/Button";
import ProcessForm from "./ui/components/process/form";
import Result from "./ui/components/Result/Result";

import { useRef, useState, useEffect } from "react";

export default function Home() {
  type ProcessResult = {
    Path: string[];
    Links: number;
    Duration: number;
    Degrees: number;
  };

  const resultSectionRef = useRef<HTMLDivElement>(null);
  const [showLoading, setShowLoading] = useState(false);
  const [isProcess, setIsProcess] = useState(false);
  const [ProcessResult, setProcessResult] = useState({});
  const [Method, setMethod] = useState("");

  const handleMethod = (method : string) => {
    setMethod(method);
    console.log(Method);
  }

  useEffect(() => {
    if (isProcess && resultSectionRef.current) {
      resultSectionRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [isProcess]);

  const handleLoading = (load : boolean) => {
    setShowLoading(load);
    setIsProcess(true);
  };

  const handleProcess = (result : ProcessResult) => {
    setProcessResult(result);
    setShowLoading(false);
  }

  return (
    <main>
      <section className="min-h-screen relative bg-gradient-to-b from-[#5356FF] via-[#378CE7] to-[#59D5E0]">
        <div className="absolute z-10 flex justify-normal px-1 pt-2 md:px-16 md:py-4">
          <Button href="/about">About</Button>
          <Button href="/how-to-use">How to use</Button>
        </div>
        <div className="min-h-screen flex flex-col md:pt-[20px] items-center relative">
          <div className="md:pl-[150px]">
            <Image
              src="/racer.png"
              alt="RACER"
              width={300}
              height={300}
              className="absolute z-0 hidden md:block"
            />
            <Image
              src="/racerMobile.png"
              alt="racerMobile"
              width={200}
              height={200}
              className="md:hidden block"
            />
          </div>
          <div className="flex flex-col relative z-10">
            <h2 className="md:text-3xl mb-[-10px] pl-4 drop-shadow-[0_1.2px_1.2px_rgba(0,0,0,0.8)] md:pt-[178px]">
              グラフレーシング
            </h2>
            <h1 className={`${racing.className} md:text-8xl drop-shadow-[0_1.2px_1.2px_rgba(0,0,0,0.8)] font-outline-2`}>
              Pathfinder Racer
            </h1>
            <h2 className={`${kotta.className} md:text-3xl mt-[-5px] ml-[-5px] font-bold drop-shadow-[0_1.2px_1.2px_rgba(0,0,0,0.8)]`}>
              Exploring the Wikipedia Universe
            </h2>
          </div>
          <ProcessForm onResult={handleProcess} showLoading={handleLoading} onMethod={handleMethod} />
        </div>
      </section>
      <div ref={resultSectionRef}>
        {isProcess && <Result Result={ProcessResult} showLoading={showLoading} method={Method}/>}
      </div>
    </main>
  );
}