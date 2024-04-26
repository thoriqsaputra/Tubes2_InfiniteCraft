export default function RadioMethod({handleOptionChange, selectedOption}) {
    return (
        <div className="flex justify-center py-5 ">
            <div className="flex items-center bg-white rounded-xl shadow-sm">
                <label className="cursor-pointer">
                    <input
                        type="radio"
                        className="peer sr-only"
                        name="solution"
                        onChange={() => handleOptionChange("all")}
                        checked={selectedOption === "all"}
                    />
                    <div className={`w-50 max-w-xl bg-white px-3 rounded-xl py-2 text-gray-600 ring-4 ring-transparent transition-all hover:shadow cursor-pointer ${selectedOption === "all" ? "peer-checked:text-white peer-checked:bg-amber-500 peer-checked:ring-offset-5" : ""}`}>
                        <p>All Solution</p>
                    </div>
                </label>
                <label className="cursor-pointer">
                    <input
                        type="radio"
                        className="peer sr-only"
                        name="solution"
                        onChange={() => handleOptionChange("one")}
                        checked={selectedOption === "one"}
                    />
                    <div className={`w-50 max-w-xl bg-white px-3 rounded-xl py-2 text-gray-600 ring-4 ring-transparent transition-all hover:shadow cursor-pointer ${selectedOption === "one" ? "peer-checked:text-white peer-checked:bg-blue-500 peer-checked:ring-offset-5" : ""}`}>
                        <p>One Solution</p>
                    </div>
                </label>
            </div>
        </div>
    )
};