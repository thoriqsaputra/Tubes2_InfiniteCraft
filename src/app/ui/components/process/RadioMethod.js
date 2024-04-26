const RadioMethod = ({handleOptionChange, selectedOption}) => {
  return (
    <div className="mx-auto max-w-6xl px-12">
      <div className="text-center mb-5">
        <h2 className="text-2xl font-bold text-[#F6F5F2] mb-1">Choose Method</h2>
        <div className={`h-1 w-[150px] bg-gradient-to-r rounded-sm ${selectedOption === "IDS" ? "from-amber-600 from bg-amber-400" : "from-blue-600 to-blue-400"} mx-auto`}></div>
      </div>
      <div className="flex flex-wrap gap-8 justify-center">
        <label className="cursor-pointer">
          <input
            type="radio"
            className="peer sr-only"
            name="method"
            onChange={() => handleOptionChange("IDS")}
            checked={selectedOption === "IDS"}
          />
          <div className={`w-50 max-w-xl rounded-md bg-white px-4 shadow-xl py-2 text-gray-600 ring-4 ring-transparent transition-all hover:shadow cursor-pointer ${selectedOption === "IDS" ? "peer-checked:text-amber-400 peer-checked:ring-amber-500 peer-checked:ring-offset-5" : ""}`}>
            <p>IDS</p>
          </div>
        </label>
        <label className="cursor-pointer">
          <input
            type="radio"
            className="peer sr-only"
            name="method"
            onChange={() => handleOptionChange("BFS")}
            checked={selectedOption === "BFS"}
          />
          <div className={`w-50 max-w-xl rounded-md bg-white px-4 shadow-xl py-2 text-gray-600 ring-4 ring-transparent transition-all hover:shadow cursor-pointer ${selectedOption === "BFS" ? "peer-checked:text-blue-400 peer-checked:ring-blue-500 peer-checked:ring-offset-5" : ""}`}>
            <p>BFS</p>
          </div>
        </label>
      </div>
    </div>
  );
};

export default RadioMethod;
