import { useState } from "react";

function Header() {
  return (
    <div className="w-full h-full justify-start pl-3 pt-3 fixed top-0 flex">
        <h1 className="text-white text-4xl font-bold font-sans">Chatbot</h1>
        <h1 className="text-[#DBEBC0] text-xl font-medium pl-3 pt-3">GPT-4o</h1>
    </div>
  );
}

export default Header;