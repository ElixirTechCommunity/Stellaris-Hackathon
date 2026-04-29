interface ConfidenceBarProps {
  confidence: number;
}

export function ConfidenceBar({ confidence }: ConfidenceBarProps) {
  const percentage = confidence * 100;
  
  return (
    <div className="w-full">
      <div className="h-1.5 bg-pink-100 rounded-none overflow-hidden">
        <div
          className="h-full bg-gradient-to-r from-pink-300 to-pink-500 transition-all duration-300"
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  );
}
