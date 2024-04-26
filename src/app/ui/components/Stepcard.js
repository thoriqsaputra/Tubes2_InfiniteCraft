const StepCard = ({ stepNumber, title }) => {
    return (
      <div className="bg-white rounded-lg shadow-md px-6 py-4 flex flex-col items-center">
        <h2 className="text-2xl text-black font-semibold mb-2">{stepNumber}</h2>
        <p className="text-base text-gray-700">{title}</p>
      </div>
    );
  };
  
  export default StepCard;