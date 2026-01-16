// Parse Tree Visualizer - Main JavaScript
// Handles Wails bindings, D3.js visualization, and step-by-step animation

import { ParseString, ParseStepByStep, ValidateGrammar, GetDefaultGrammar, GetTokens } from '../wailsjs/go/main/App';

// State
let currentTree = null;
let currentSteps = [];
let currentStepIndex = 0;
let isPlaying = false;
let playInterval = null;
let animationSpeed = 500;
let svg = null;
let g = null;
let zoom = null;

// DOM Elements
const grammarInput = document.getElementById('grammarInput');
const inputString = document.getElementById('inputString');
const parseBtn = document.getElementById('parseBtn');
const stepBtn = document.getElementById('stepBtn');
const loadDefaultBtn = document.getElementById('loadDefaultBtn');
const tokenDisplay = document.getElementById('tokenDisplay');
const grammarStatus = document.getElementById('grammarStatus');
const treeContainer = document.getElementById('treeContainer');
const treeSvg = document.getElementById('treeSvg');
const errorDisplay = document.getElementById('errorDisplay');
const stepControls = document.getElementById('stepControls');
const stepCounter = document.getElementById('stepCounter');
const stepDescription = document.getElementById('stepDescription');
const prevStepBtn = document.getElementById('prevStepBtn');
const nextStepBtn = document.getElementById('nextStepBtn');
const playPauseBtn = document.getElementById('playPauseBtn');
const playIcon = document.getElementById('playIcon');
const pauseIcon = document.getElementById('pauseIcon');
const resetStepBtn = document.getElementById('resetStepBtn');
const speedSlider = document.getElementById('speedSlider');
const speedValue = document.getElementById('speedValue');
const zoomInBtn = document.getElementById('zoomInBtn');
const zoomOutBtn = document.getElementById('zoomOutBtn');
const resetZoomBtn = document.getElementById('resetZoomBtn');
const helpBtn = document.getElementById('helpBtn');
const helpModal = document.getElementById('helpModal');

// Initialize
document.addEventListener('DOMContentLoaded', init);

function init() {
    setupEventListeners();
    setupD3();
    loadDefaultGrammar();
}

function setupEventListeners() {
    parseBtn.addEventListener('click', handleParse);
    stepBtn.addEventListener('click', handleStepMode);
    loadDefaultBtn.addEventListener('click', loadDefaultGrammar);
    inputString.addEventListener('input', handleInputChange);
    grammarInput.addEventListener('input', handleGrammarChange);
    
    // Step controls
    prevStepBtn.addEventListener('click', handlePrevStep);
    nextStepBtn.addEventListener('click', handleNextStep);
    playPauseBtn.addEventListener('click', handlePlayPause);
    resetStepBtn.addEventListener('click', handleResetSteps);
    speedSlider.addEventListener('input', handleSpeedChange);
    
    // Zoom controls
    zoomInBtn.addEventListener('click', () => zoomBy(1.3));
    zoomOutBtn.addEventListener('click', () => zoomBy(0.7));
    resetZoomBtn.addEventListener('click', resetZoom);
    
    // Help modal
    helpBtn.addEventListener('click', () => helpModal.classList.remove('hidden'));
    helpModal.querySelector('.modal-backdrop').addEventListener('click', () => helpModal.classList.add('hidden'));
    helpModal.querySelector('.modal-close').addEventListener('click', () => helpModal.classList.add('hidden'));
}

function setupD3() {
    const width = treeContainer.clientWidth;
    const height = treeContainer.clientHeight;
    
    svg = d3.select('#treeSvg')
        .attr('width', width)
        .attr('height', height);
    
    // Create zoom behavior
    zoom = d3.zoom()
        .scaleExtent([0.1, 4])
        .on('zoom', (event) => {
            g.attr('transform', event.transform);
        });
    
    svg.call(zoom);
    
    // Create main group for tree
    g = svg.append('g')
        .attr('transform', `translate(${width / 2}, 50)`);
    
    // Handle window resize
    window.addEventListener('resize', () => {
        svg.attr('width', treeContainer.clientWidth)
           .attr('height', treeContainer.clientHeight);
    });
}

async function loadDefaultGrammar() {
    try {
        const grammar = await GetDefaultGrammar();
        grammarInput.value = grammar;
        validateGrammar();
    } catch (err) {
        console.error('Failed to load default grammar:', err);
    }
}

async function handleGrammarChange() {
    await validateGrammar();
}

async function validateGrammar() {
    const grammarText = grammarInput.value.trim();
    if (!grammarText) {
        grammarStatus.textContent = '';
        grammarStatus.className = 'status-message';
        return;
    }
    
    try {
        const result = await ValidateGrammar(grammarText);
        if (result.valid) {
            grammarStatus.textContent = '✓ Grammar is valid';
            grammarStatus.className = 'status-message success';
        } else {
            grammarStatus.textContent = '✗ ' + result.errors.join('; ');
            grammarStatus.className = 'status-message error';
        }
        
        if (result.warnings && result.warnings.length > 0) {
            grammarStatus.textContent += ' ⚠ ' + result.warnings.join('; ');
            if (result.valid) {
                grammarStatus.className = 'status-message warning';
            }
        }
    } catch (err) {
        grammarStatus.textContent = '✗ ' + err;
        grammarStatus.className = 'status-message error';
    }
}

async function handleInputChange() {
    const input = inputString.value.trim();
    if (!input) {
        tokenDisplay.innerHTML = '';
        return;
    }
    
    try {
        const tokens = await GetTokens(input);
        renderTokens(tokens);
    } catch (err) {
        console.error('Failed to tokenize:', err);
    }
}

function renderTokens(tokens) {
    tokenDisplay.innerHTML = tokens
        .filter(t => t.type !== 'EOF')
        .map(t => {
            let className = 'token';
            if (t.type === 'NUMBER') className += ' number';
            else if (['PLUS', 'MINUS', 'MULT', 'DIV'].includes(t.type)) className += ' operator';
            else if (['LPAREN', 'RPAREN'].includes(t.type)) className += ' paren';
            return `<span class="${className}">${t.value}</span>`;
        })
        .join('');
}

async function handleParse() {
    const grammarText = grammarInput.value.trim();
    const input = inputString.value.trim();
    
    if (!grammarText) {
        showError('Please enter a grammar');
        return;
    }
    
    if (!input) {
        showError('Please enter an input string');
        return;
    }
    
    hideError();
    hideStepControls();
    
    try {
        const result = await ParseString(grammarText, input);
        
        if (result.success) {
            currentTree = result.tree;
            renderTree(currentTree);
            hideEmptyState();
        } else {
            showError(result.error);
            clearTree();
        }
    } catch (err) {
        showError('Parse error: ' + err);
        clearTree();
    }
}

async function handleStepMode() {
    const grammarText = grammarInput.value.trim();
    const input = inputString.value.trim();
    
    if (!grammarText) {
        showError('Please enter a grammar');
        return;
    }
    
    if (!input) {
        showError('Please enter an input string');
        return;
    }
    
    hideError();
    
    try {
        const result = await ParseStepByStep(grammarText, input);
        
        if (result.success) {
            currentTree = result.tree;
            currentSteps = result.steps;
            currentStepIndex = 0;
            isPlaying = false;
            
            showStepControls();
            hideEmptyState();
            clearTree();
            updateStepUI();
        } else {
            showError(result.error);
            hideStepControls();
        }
    } catch (err) {
        showError('Parse error: ' + err);
        hideStepControls();
    }
}

function handlePrevStep() {
    if (currentStepIndex > 0) {
        currentStepIndex--;
        renderPartialTree();
        updateStepUI();
    }
}

function handleNextStep() {
    if (currentStepIndex < currentSteps.length) {
        currentStepIndex++;
        renderPartialTree();
        updateStepUI();
    }
}

function handlePlayPause() {
    if (isPlaying) {
        stopAutoPlay();
    } else {
        startAutoPlay();
    }
}

function startAutoPlay() {
    isPlaying = true;
    playIcon.classList.add('hidden');
    pauseIcon.classList.remove('hidden');
    
    playInterval = setInterval(() => {
        if (currentStepIndex < currentSteps.length) {
            currentStepIndex++;
            renderPartialTree();
            updateStepUI();
        } else {
            stopAutoPlay();
        }
    }, animationSpeed);
}

function stopAutoPlay() {
    isPlaying = false;
    playIcon.classList.remove('hidden');
    pauseIcon.classList.add('hidden');
    
    if (playInterval) {
        clearInterval(playInterval);
        playInterval = null;
    }
}

function handleResetSteps() {
    stopAutoPlay();
    currentStepIndex = 0;
    clearTree();
    updateStepUI();
}

function handleSpeedChange() {
    animationSpeed = parseInt(speedSlider.value);
    speedValue.textContent = animationSpeed + 'ms';
    
    // Restart auto-play with new speed if playing
    if (isPlaying) {
        stopAutoPlay();
        startAutoPlay();
    }
}

function updateStepUI() {
    stepCounter.textContent = `Step ${currentStepIndex} of ${currentSteps.length}`;
    
    if (currentStepIndex > 0 && currentStepIndex <= currentSteps.length) {
        const step = currentSteps[currentStepIndex - 1];
        stepDescription.textContent = step.description;
    } else {
        stepDescription.textContent = 'Ready to start';
    }
    
    prevStepBtn.disabled = currentStepIndex === 0;
    nextStepBtn.disabled = currentStepIndex >= currentSteps.length;
}

function renderPartialTree() {
    if (currentStepIndex === 0) {
        clearTree();
        return;
    }
    
    // Build tree up to current step
    const nodeMap = new Map();
    const rootNodes = [];
    
    for (let i = 0; i < currentStepIndex; i++) {
        const step = currentSteps[i];
        const node = {
            id: step.nodeId,
            label: step.description.includes('terminal') ? 
                step.description.match(/'([^']+)'/)?.[1] || 'ε' :
                step.description.match(/<([^>]+)>/)?.[1] || step.description,
            isTerminal: step.description.includes('terminal') || step.description.includes('epsilon'),
            children: [],
            isNew: i === currentStepIndex - 1
        };
        
        nodeMap.set(step.nodeId, node);
        
        if (step.parentId === -1 || step.parentId === 0) {
            rootNodes.push(node);
        } else if (nodeMap.has(step.parentId)) {
            nodeMap.get(step.parentId).children.push(node);
        } else {
            rootNodes.push(node);
        }
    }
    
    // Create a proper tree structure for D3
    if (rootNodes.length > 0) {
        const tree = rootNodes[0];
        renderTree(tree, true);
    }
}

function renderTree(data, animate = false) {
    if (!data) return;
    
    // Clear existing tree
    g.selectAll('*').remove();
    
    // Create D3 hierarchy
    const root = d3.hierarchy(data, d => d.children);
    
    // Calculate tree layout
    const treeLayout = d3.tree()
        .nodeSize([80, 100])
        .separation((a, b) => a.parent === b.parent ? 1 : 1.5);
    
    treeLayout(root);
    
    // Draw links first (so they appear behind nodes)
    const links = g.selectAll('.link')
        .data(root.links())
        .enter()
        .append('path')
        .attr('class', d => 'link' + (animate && d.target.data.isNew ? ' new' : ''))
        .attr('d', d3.linkVertical()
            .x(d => d.x)
            .y(d => d.y));
    
    // Draw nodes
    const nodes = g.selectAll('.node')
        .data(root.descendants())
        .enter()
        .append('g')
        .attr('class', 'node')
        .attr('transform', d => `translate(${d.x}, ${d.y})`);
    
    // Node circles
    nodes.append('circle')
        .attr('r', 25)
        .attr('class', d => {
            let classes = '';
            if (d.data.isTerminal) classes += 'terminal ';
            if (animate && d.data.isNew) classes += 'new';
            return classes.trim();
        });
    
    // Node labels
    nodes.append('text')
        .attr('class', d => d.data.isTerminal ? 'terminal' : '')
        .attr('dy', '0.35em')
        .text(d => {
            const label = d.data.label || d.data.Label;
            // Truncate long labels
            return label && label.length > 8 ? label.substring(0, 6) + '...' : label;
        });
    
    // Center the tree
    centerTree(root);
}

function centerTree(root) {
    if (!root) return;
    
    const bounds = g.node().getBBox();
    const width = treeContainer.clientWidth;
    const height = treeContainer.clientHeight;
    
    const scale = Math.min(
        (width - 100) / bounds.width,
        (height - 100) / bounds.height,
        1
    );
    
    const x = width / 2 - (bounds.x + bounds.width / 2) * scale;
    const y = 50;
    
    svg.transition()
        .duration(300)
        .call(zoom.transform, d3.zoomIdentity.translate(x, y).scale(scale));
}

function clearTree() {
    g.selectAll('*').remove();
}

function zoomBy(factor) {
    svg.transition()
        .duration(300)
        .call(zoom.scaleBy, factor);
}

function resetZoom() {
    if (currentTree) {
        const root = d3.hierarchy(currentTree, d => d.children);
        centerTree(root);
    } else {
        svg.transition()
            .duration(300)
            .call(zoom.transform, d3.zoomIdentity.translate(treeContainer.clientWidth / 2, 50));
    }
}

function showError(message) {
    errorDisplay.textContent = message;
    errorDisplay.classList.remove('hidden');
}

function hideError() {
    errorDisplay.classList.add('hidden');
}

function showStepControls() {
    stepControls.classList.remove('hidden');
}

function hideStepControls() {
    stepControls.classList.add('hidden');
    stopAutoPlay();
}

function hideEmptyState() {
    const emptyState = treeContainer.querySelector('.empty-state');
    if (emptyState) {
        emptyState.style.display = 'none';
    }
}
