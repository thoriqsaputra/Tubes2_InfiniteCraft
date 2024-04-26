import Button from '../ui/components/Button';
import StepCard from '../ui/components/StepCard';

export default function Page() {
    return (
        <div className="min-h-screen bg-gradient-to-br from-purple-500 to-blue-500 ">
            <div className="flex justify-normal px-1 pt-2 md:px-16 md:py-4 mb-10">
              <Button href="/about">About Us</Button>
              <Button href="/">Home</Button>
            </div>
    
          <div className="max-w-4xl mx-auto px-4 pb-8 flex flex-col items-center justify-center">
            <h1 className="text-4xl font-bold text-white text-center mb-8">
              How to Use Our Tool
            </h1>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <StepCard stepNumber="1" title="Choose Your Language" />
              <StepCard stepNumber="2" title="Select Pathfinding Method" />
              <StepCard stepNumber="3" title="Enter Start & Destination Links" />
              <StepCard stepNumber="4" title="Click 'Start Race' Button" />
              <StepCard stepNumber="5" title="View Your Calculated Path" />
            </div>
          </div>
        </div>
      );
  }
  