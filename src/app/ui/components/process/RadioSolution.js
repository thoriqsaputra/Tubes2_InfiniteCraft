import { useState } from 'react';
import { motion } from 'framer-motion';

export default function RadioMethod({ handleOptionChange, selectedOption }) {
    const [hoveredOption, setHoveredOption] = useState(null);

    const variants = {
        checked: {
            scale: 1.05,
            backgroundColor: hoveredOption || selectedOption === "all" ? "#FF6B6B" : "#074173",
            color: "#FFFFFF",
            transition: {
                duration: 0.2,
                type: "spring",
                stiffness: 200,
            }
        },
        unchecked: {
            scale: 1,
            backgroundColor: hoveredOption || selectedOption === "all" ? "#F3F4F6" : "#E0E7FF",
            color: "#4B5563",
            transition: {
                duration: 0.2,
                type: "spring",
                stiffness: 200,
            }
        }
    };

    return (
        <div className="flex justify-center py-5">
            <div className="flex items-center space-x-4">
                <motion.label
                    className="cursor-pointer"
                    whileHover={{ scale: 1.05 }}
                >
                    <input
                        type="radio"
                        className="peer sr-only"
                        name="solution"
                        onChange={() => handleOptionChange("all")}
                        checked={selectedOption === "all"}
                        onMouseEnter={() => setHoveredOption("all")}
                        onMouseLeave={() => setHoveredOption(null)}
                    />
                    <motion.div
                        className={`w-36 bg-white px-3 rounded-xl py-2 text-gray-600 ring-4 ring-transparent cursor-pointer flex justify-center items-center`}
                        variants={variants}
                        animate={selectedOption === "all" ? "checked" : "unchecked"}
                    >
                        <p>All Solution</p>
                    </motion.div>
                </motion.label>
                <motion.label
                    className="cursor-pointer"
                    whileHover={{ scale: 1.05 }}
                >
                    <input
                        type="radio"
                        className="peer sr-only"
                        name="solution"
                        onChange={() => handleOptionChange("one")}
                        checked={selectedOption === "one"}
                        onMouseEnter={() => setHoveredOption("one")}
                        onMouseLeave={() => setHoveredOption(null)}
                    />
                    <motion.div
                        className={`w-36 bg-white px-3 rounded-xl py-2 text-gray-600 ring-4 ring-transparent cursor-pointer flex justify-center items-center`}
                        variants={variants}
                        animate={selectedOption === "one" ? "checked" : "unchecked"}
                    >
                        <p>One Solution</p>
                    </motion.div>
                </motion.label>
            </div>
        </div>
    );
}