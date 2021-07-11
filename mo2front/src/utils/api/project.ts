import { ParseQuery } from "./queryParser";
import axios from "axios";
import { Project } from "../../models";



export async function ListProject(q: { page: number, pageSize: number, tags?: string[], uid?: string }) {
    return (await axios.get<Project[]>('/api/project' + ParseQuery(q))).data
}