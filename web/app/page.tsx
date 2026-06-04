'use client';

import { useState } from 'react';

export default function RagDashboard() {
  const [ingestStatus, setIngestStatus] = useState<string>('Ready to ingest documents.');
  const [isDragging, setIsDragging] = useState(false);
  
  const [query, setQuery] = useState('');
  // NEW: State to hold our comma-separated metadata filters (e.g., "author: Thorne, year: 2084")
  const [filterInput, setFilterInput] = useState(''); 
  const [isQuerying, setIsQuerying] = useState(false);
  
  const [result, setResult] = useState<{
    answer: string, 
    contexts_used: number, 
    search_pool_size: number,
    contexts?: string[] 
  } | null>(null);

  const processFile = async (file: File) => {
    setIngestStatus(`Reading ${file.name}...`);
    
    const reader = new FileReader();
    reader.onload = async (e) => {
      const text = e.target?.result;
      if (typeof text === 'string') {
        setIngestStatus('Generating embeddings and indexing in AuraVector-Go...');
        try {
          const res = await fetch('http://localhost:8080/rag/ingest', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ text }), // Note: Real-world apps would also pass metadata tags here!
          });
          const data = await res.json();
          setIngestStatus(`✅ Success: Generated and indexed ${data.chunks_processed} semantic chunks.`);
        } catch (error) {
          setIngestStatus(`❌ Connection Error: Is the Go engine running?`);
        }
      }
    };
    reader.onerror = () => setIngestStatus("❌ Failed to read file.");
    reader.readAsText(file);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      processFile(e.dataTransfer.files[0]);
    }
  };

  const executeQuery = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    
    setIsQuerying(true);
    setResult(null);

    // Parse the user's string input "key: value, key2: value2" into a JSON object
    const parsedFilters: Record<string, string> = {};
    if (filterInput.trim()) {
        const pairs = filterInput.split(',');
        pairs.forEach(pair => {
            const [key, value] = pair.split(':');
            if (key && value) {
                parsedFilters[key.trim()] = value.trim();
            }
        });
    }
    
    try {
      const res = await fetch('http://localhost:8080/rag/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        // NEW: Send the parsed filters to our Phase F Go Backend
        body: JSON.stringify({ query, k: 3, filters: parsedFilters }), 
      });
      const data = await res.json();
      setResult(data);
    } catch (error) {
      console.error("Query failed", error);
    } finally {
      setIsQuerying(false);
    }
  };

  return (
    <main className="min-h-screen bg-gray-950 text-gray-100 p-8 font-sans selection:bg-purple-500/30">
      <header className="max-w-5xl mx-auto mb-10 border-b border-gray-800 pb-6 flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
            AuraVector-Go Intelligence
          </h1>
          <p className="text-gray-400 mt-2">Private Document Semantic Search Engine (RAG)</p>
        </div>
      </header>

      <div className="max-w-5xl mx-auto grid grid-cols-1 md:grid-cols-3 gap-8">
        {/* Document Ingestion Zone */}
        <div className="flex flex-col gap-4">
          <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-2xl">
            <h2 className="text-lg font-semibold text-blue-400 mb-4 flex items-center gap-2">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path></svg>
              Knowledge Ingestion
            </h2>
            
            <div 
              onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
              onDragLeave={() => setIsDragging(false)}
              onDrop={handleDrop}
              className={`border-2 border-dashed rounded-lg p-8 text-center transition-all duration-200 ease-in-out ${isDragging ? 'border-purple-500 bg-purple-500/10' : 'border-gray-700 bg-gray-800/50 hover:border-gray-500'}`}
            >
              <svg className="w-10 h-10 mx-auto text-gray-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
              <p className="text-sm text-gray-300 font-medium">Drag & Drop Documents Here</p>
              <p className="text-xs text-gray-500 mt-2">Supports .txt, .md, .csv</p>
            </div>
            
            <div className="mt-4 p-3 bg-black/40 rounded border border-gray-800">
              <span className="block text-xs uppercase tracking-wider text-gray-500 mb-1">Pipeline Status</span>
              <span className="text-sm text-blue-300 font-mono">{ingestStatus}</span>
            </div>
          </div>
        </div>

        {/* Semantic Query Workspace */}
        <div className="md:col-span-2 flex flex-col gap-4">
          <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-2xl min-h-[500px] flex flex-col">
            <h2 className="text-lg font-semibold text-purple-400 mb-4 flex items-center gap-2">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14 10l-2 1m0 0l-2-1m2 1v2.5M20 7l-2 1m2-1l-2-1m2 1v2.5M14 4l-2-1-2 1M4 7l2-1M4 7l2 1M4 7v2.5M12 21l-2-1m2 1l2-1m-2 1v-2.5M6 18l-2-1v-2.5M18 18l2-1v-2.5"></path></svg>
              Semantic RAG Workspace
            </h2>

            <form onSubmit={executeQuery} className="flex flex-col gap-3 mb-6">
              
              {/* NEW: Metadata Filter Input */}
              <div className="flex gap-2 items-center bg-gray-950 border border-gray-800 rounded-lg p-2 focus-within:border-gray-600 transition-colors">
                 <svg className="w-4 h-4 text-gray-500 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"></path></svg>
                 <input
                  type="text"
                  placeholder="Optional Filters (e.g., author: Thorne, category: science)"
                  value={filterInput}
                  onChange={(e) => setFilterInput(e.target.value)}
                  className="flex-1 bg-transparent text-xs text-gray-300 focus:outline-none placeholder:text-gray-600"
                  disabled={isQuerying}
                 />
              </div>

              <div className="flex gap-3">
                <input
                  type="text"
                  placeholder="Ask a question about your documents..."
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  className="flex-1 bg-gray-950 text-gray-100 border border-gray-700 p-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all placeholder:text-gray-600"
                  disabled={isQuerying}
                />
                <button 
                  type="submit" 
                  disabled={isQuerying || !query.trim()}
                  className="bg-purple-600 hover:bg-purple-500 disabled:bg-gray-700 disabled:text-gray-500 text-white px-6 py-3 rounded-lg font-semibold transition-colors flex items-center gap-2"
                >
                  {isQuerying ? (
                    <span className="animate-pulse">Synthesizing...</span>
                  ) : (
                    <>Ask LLM <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg></>
                  )}
                </button>
              </div>
            </form>

            <div className="flex-1 bg-black/40 border border-gray-800 rounded-lg p-6 overflow-y-auto">
              {result ? (
                <div className="animate-in fade-in slide-in-from-bottom-4 duration-500 flex flex-col h-full">
                  
                  {/* Metadata Header */}
                  <div className="flex items-center gap-4 text-xs uppercase tracking-wider text-gray-500 mb-4 border-b border-gray-800 pb-3 shrink-0">
                    <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-green-500"></span> Contexts Used: {result.contexts_used}</span>
                    <span className="flex items-center gap-1"><span className="w-2 h-2 rounded-full bg-blue-500"></span> Total Vector Pool: {result.search_pool_size}</span>
                  </div>
                  
                  {/* LLM Generated Answer */}
                  <div className="prose prose-invert max-w-none mb-6">
                    <p className="text-gray-200 leading-relaxed whitespace-pre-wrap">
                      {result.answer}
                    </p>
                  </div>

                  {/* Context Reference Viewer */}
                  {result.contexts && result.contexts.length > 0 && (
                    <div className="mt-auto border-t border-gray-800 pt-5 shrink-0">
                      <h3 className="text-xs font-bold text-gray-500 uppercase tracking-wider mb-3 flex items-center gap-2">
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
                        Source Context Rendered
                      </h3>
                      <div className="flex flex-col gap-3">
                        {result.contexts.map((ctx, index) => (
                          <div key={index} className="bg-gray-950 p-4 rounded-lg border border-gray-700/50 shadow-inner">
                            <span className="text-blue-400 font-mono text-xs font-semibold mb-2 block">
                              [Source Chunk {index + 1}]
                            </span>
                            <p className="text-sm text-gray-400 italic font-serif leading-relaxed">
                              "{ctx}"
                            </p>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                </div>
              ) : (
                <div className="h-full flex flex-col items-center justify-center text-gray-600">
                  <svg className="w-12 h-12 mb-3 opacity-20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
                  <p>Awaiting your query...</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}