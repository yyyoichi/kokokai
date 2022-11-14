import {
    Node,
    Edge,
    Data as NodeData,
} from "vis-network/standalone/umd/vis-network.min.js";

export default class NetworkNode {
    private edges: Edge[] = [];
    private nodeHash: string[] = [];
    constructor() {}
    /**
     *
     * @param label 頂点
     * @returns 追加した番号
     */
    addNode(label: string) {
        let l = this.nodeHash.indexOf(label);
        if (l === -1) {
            this.nodeHash.push(label);
            l = this.nodeHash.length;
        }
        return l;
    }
    /**
     * 頂点と頂点のペアをセット
     * @param label1 繋げたい頂点
     * @param label2 繋げたい頂点
     */
    addEdge(label1: string, label2: string) {
        const l1 = this.nodeHash.indexOf(label1);
        const l2 = this.nodeHash.indexOf(label2);
        this.edges.push({ from: l1, to: l2 });
    }
    getNodes(): Node[] {
        return this.nodeHash.map((label, i) => {
            return { id: i + 1, label };
        });
    }
    getEdges(): Edge[] {
        return this.edges;
    }
    getNodeData() {
        return {
            nodes: this.getNodes(),
            edges: this.getEdges(),
        };
    }
}
