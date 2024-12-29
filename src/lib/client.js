import axios from "axios";
import {HOST} from "../routes/routes.js";

export const apiClient = axios.create({
    baseURL: HOST,
})
