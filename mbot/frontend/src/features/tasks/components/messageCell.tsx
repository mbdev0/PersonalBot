const TX_HASH_REGEX = /[1-9A-HJ-NP-Za-km-z]{87,88}/g;
const SOLSCAN_TX_URL = 'https://solscan.io/tx/';

function parseParts(message: string): Array<{ type: 'text' | 'tx'; value: string }> {
  const parts: Array<{ type: 'text' | 'tx'; value: string }> = [];
  let lastIndex = 0;

  for (const match of message.matchAll(TX_HASH_REGEX)) {
    if (match.index > lastIndex) {
      parts.push({ type: 'text', value: message.slice(lastIndex, match.index) });
    }
    parts.push({ type: 'tx', value: match[0] });
    lastIndex = match.index + match[0].length;
  }

  if (lastIndex < message.length) {
    parts.push({ type: 'text', value: message.slice(lastIndex) });
  }

  return parts;
}

export function MessageCell({ message }: { message: string }) {
  if (!message) {
    return <div className="font-mono text-xs text-foreground/30 italic">...</div>;
  }

  const parts = parseParts(message);

  return (
    <div
      title={message}
      className="font-mono text-xs bg-foreground/5 border border-foreground/10
                 rounded px-2 py-1 truncate cursor-default "
    >
      {parts.map((part, i) =>
        part.type === 'tx' ? (
          <a
            key={i}
            href={`${SOLSCAN_TX_URL}${part.value}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-400 hover:text-blue-300 underline underline-offset-2"
          >
            solscan ↗
          </a>
        ) : (
          <span key={i}>{part.value}</span>
        )
      )}
    </div>
  );
}
