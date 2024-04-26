'use client';

import Link from 'next/link';

const Button = ({ href, children,  }) => {
return (
    <Link href={href}
            className=" text-white font-bold mx-1 py-2 px-4 inline-block transition duration-500 transform hover:translate-y-2 rounded hover:bg-[#1679AB] hover:text-[#DFF5FF]">
            {children}
    </Link>
);
};

export default Button;
