const RadioLaguage = ({handleOptionChange, selectedOption}) => {
  return (
    <div className="mx-auto max-w-6xl px-12">
      <div className="text-center mb-5">
        <h2 className="text-2xl font-bold text-[#F6F5F2] mb-1">Choose Language</h2>
        <div className={`h-1 w-[200px] bg-gradient-to-r rounded-sm ${selectedOption === "id" ? "from-red-400 from bg-red-600" : "from-green-400 to-green-600"} mx-auto`}></div>
      </div>
      <div className="flex flex-wrap gap-8 justify-center">
        <label className="cursor-pointer">
          <input
            type="radio"
            className="peer sr-only"
            name="language"
            onChange={() => handleOptionChange("id")}
            checked={selectedOption === "id"}
          />
          <div className={`w-50 max-w-xl rounded-md bg-white px-4 shadow-xl py-2 text-gray-600 ring-4 ring-transparent transition-all hover:shadow cursor-pointer ${selectedOption === "id" ? "peer-checked:text-red-500 peer-checked:ring-red-600 peer-checked:ring-offset-5" : ""}`}>
            <p>Indonesian</p>
          </div>
        </label>
        <label className="cursor-pointer">
          <input
            type="radio"
            className="peer sr-only"
            name="language"
            onChange={() => handleOptionChange("en")}
            checked={selectedOption === "en"}
          />
          <div className={`w-50 max-w-xl rounded-md bg-white px-4 shadow-xl py-2 text-gray-600 ring-4 ring-transparent transition-all hover:shadow cursor-pointer ${selectedOption === "en" ? "peer-checked:text-green-400 peer-checked:ring-green-500 peer-checked:ring-offset-5" : ""}`}>
            <p>English</p>
          </div>
        </label>
      </div>
    </div>
  );
};

export default RadioLaguage;
