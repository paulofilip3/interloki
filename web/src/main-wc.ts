import { defineCustomElement } from 'vue'
import InterlokiViewer from './InterlokiViewer.ce.vue'

const InterlokiViewerElement = defineCustomElement(InterlokiViewer)
customElements.define('interloki-viewer', InterlokiViewerElement)
