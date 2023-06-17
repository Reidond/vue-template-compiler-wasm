import compiler from "vue-template-compiler";
import { outdent } from "outdent";

// Read input from stdin
function readInput() {
  const chunkSize = 1024;
  const inputChunks = [];
  let totalBytes = 0;

  // Read all the available bytes
  while (1) {
    const buffer = new Uint8Array(chunkSize);
    // Stdin file descriptor
    const fd = 0;
    const bytesRead = Javy.IO.readSync(fd, buffer);

    totalBytes += bytesRead;
    if (bytesRead === 0) {
      break;
    }
    inputChunks.push(buffer.subarray(0, bytesRead));
  }

  // Assemble input into a single Uint8Array
  const { finalBuffer } = inputChunks.reduce(
    (context, chunk) => {
      context.finalBuffer.set(chunk, context.bufferOffset);
      context.bufferOffset += chunk.length;
      return context;
    },
    { bufferOffset: 0, finalBuffer: new Uint8Array(totalBytes) }
  );

  return new TextDecoder().decode(finalBuffer);
}

// Write output to stdout
function writeOutput(output) {
  const encodedOutput = new TextEncoder().encode(output);
  const buffer = new Uint8Array(encodedOutput);
  // Stdout file descriptor
  const fd = 1;
  Javy.IO.writeSync(fd, buffer);
}

function getExportDefaultCode(code) {
  const regex = /export\s+default\s+{([\s\S]*?)};/;
  const match = code.match(regex);

  if (!match) {
    throw new Error("Cannot get export default contents");
  }

  const moduleCode = match[1];
  return moduleCode;
}

/**
 * createVueApp
 * @param {string} mountId
 * @param {compiler.SFCDescriptor} componentDescriptor
 * @returns {string}
 */
function createVueApp(mountId, componentDescriptor) {
  const code = getExportDefaultCode(componentDescriptor.script.content);

  const appCode = outdent`
    import { createApp } from 'https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js'
    const app = createApp({
      template: '${componentDescriptor.template.content}',
      ${code}
    });
    app.mount("#${mountId}");\n
  `;
  const styleCode = componentDescriptor.styles.map(
    (style) => `#${mountId} { ${style.content} }`
  );

  return {
    appCode,
    styleCode,
  };
}

// Read input from stdin
const input = JSON.parse(readInput());
if (!input?.sfcCode && !input?.mountId) {
  throw new Error(
    `"sfcCode" in input is empty or missing; "mountId" in input is empty or missing`
  );
}

// Call the function with the input
const sfcDescriptor = compiler.parseComponent(input.sfcCode);
const vueApp = createVueApp(input.mountId, sfcDescriptor);

// Write the result to stdout
writeOutput(JSON.stringify(vueApp));
