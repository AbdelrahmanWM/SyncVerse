import React, {  useRef, useState } from "react";

const TextAreaBasedTextEditor: React.FC = () => {
  const [value, setValue] = useState<string>("");

  const prevValueRef = useRef<string>("");
  const prevSelectionStartRef = useRef<number>(0);
  const prevSelectionEndRef = useRef<number>(0);
  const isBackspaceRef = useRef<boolean>(false);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    const { selectionStart = 0, selectionEnd = 0 } = e.currentTarget;

    prevSelectionStartRef.current = selectionStart;
    prevSelectionEndRef.current = selectionEnd;

    isBackspaceRef.current = e.key === "Backspace";
  };

  const handleInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const target = e.currentTarget;
    const newValue = target.value;
    const currSelectionStart = target.selectionStart ?? 0;

    const prevValue = prevValueRef.current;
    const prevStart = prevSelectionStartRef.current;
    const prevEnd = prevSelectionEndRef.current;
    const isBackspace = isBackspaceRef.current;

    let insertedText = "";
    let removedText = "";
    let action = "";

    const hadSelection = prevStart !== prevEnd;

    if (isBackspace && hadSelection) {
      // Deletion over selection via backspace
      removedText = prevValue.slice(prevStart, prevEnd);
      action = "deleted";
    } else if (isBackspace && !hadSelection) {
      // Normal single-char backspace
      const removedLength = prevValue.length - newValue.length;
      removedText = prevValue.slice(prevStart - removedLength, prevStart);
      action = "deleted";
    } else {
      insertedText = newValue.slice(prevStart, currSelectionStart);
      const removedLength =
        prevValue.length +
        (currSelectionStart - prevStart) -
        newValue.length;

      removedText =
        removedLength > 0
          ? prevValue.slice(prevStart, prevStart + removedLength)
          : "";
      action =
        insertedText.length > 0
          ? removedText
            ? "replaced"
            : "inserted"
          : "deleted";
    }

    console.log({
      action,
      insertedText,
      removedText,
      position: prevStart,
    });

    // Update state and ref
    prevValueRef.current = newValue;
    setValue(newValue);
  };

  return (
    <textarea
      rows={10}
      cols={50}
      value={value}
      onKeyDown={handleKeyDown}
      onInput={handleInput}
    />
  );
};

export default TextAreaBasedTextEditor;

