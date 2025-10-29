// API client for TF-Engine backend
import { logger } from '$lib/utils/logger';

const API_BASE_URL = '/api';

interface ApiResponse<T> {
	data: T;
}

interface ApiError {
	error: string;
	message?: string;
	code: number;
}

export interface Settings {
	equity: number;
	riskPct: number;
	portfolioCap: number;
	bucketCap: number;
	maxUnits: number;
}

export interface Position {
	id: number;
	ticker: string;
	entry_price: number;
	current_stop: number;
	initial_stop: number;
	shares: number;
	risk_dollars: number;
	bucket?: string;
	status: string;
	decision_id: number;
	opened_at: string;
}

export interface PositionInfo {
	ticker: string;
	entry_price: number;
	risk_dollars: number;
	status: string;
	days_held: number;
}

export interface WeekData {
	week_start: string;
	week_end: string;
	sectors: Record<string, PositionInfo[]>;
}

export interface CalendarData {
	weeks: WeekData[];
	sectors: string[];
}

export interface Candidate {
	id: number;
	ticker: string;
	date: string;
	sector?: string;
	bucket?: string;
}

class ApiClient {
	private async request<T>(
		method: string,
		endpoint: string,
		body?: unknown
	): Promise<ApiResponse<T>> {
		const url = `${API_BASE_URL}${endpoint}`;
		const startTime = performance.now();

		logger.apiRequest(method, url);

		try {
			const response = await fetch(url, {
				method,
				headers: {
					'Content-Type': 'application/json'
				},
				body: body ? JSON.stringify(body) : undefined
			});

			const duration = Math.round(performance.now() - startTime);

			if (!response.ok) {
				logger.apiResponse(method, url, response.status, duration);
				const error: ApiError = await response.json();
				throw new Error(error.message || error.error);
			}

			const data = await response.json();
			logger.apiResponse(method, url, response.status, duration);

			return data as ApiResponse<T>;
		} catch (error) {
			const duration = Math.round(performance.now() - startTime);
			logger.error(`API ${method} ${url} failed`, error);
			throw error;
		}
	}

	// Settings
	async getSettings(): Promise<Settings> {
		const response = await this.request<Settings>('GET', '/settings');
		return response.data;
	}

	// Positions
	async getPositions(): Promise<Position[]> {
		const response = await this.request<Position[]>('GET', '/positions');
		return response.data;
	}

	// Candidates
	async getCandidates(date?: string): Promise<Candidate[]> {
		const dateParam = date ? `?date=${date}` : '';
		const response = await this.request<Candidate[]>('GET', `/candidates${dateParam}`);
		return response.data;
	}

	async importCandidates(tickers: string[], date: string): Promise<{ imported: number; date: string }> {
		const response = await this.request<{ imported: number; date: string }>(
			'POST',
			'/candidates/import',
			{ tickers, date }
		);
		return response.data;
	}

	async scanCandidates(preset: string): Promise<{ count: number; tickers: string[]; date: string }> {
		const response = await this.request<{ count: number; tickers: string[]; date: string }>(
			'POST',
			'/candidates/scan',
			{ preset }
		);
		return response.data;
	}

	// Position Sizing
	async calculateSize(request: SizingRequest): Promise<SizingResult> {
		const response = await this.request<SizingResult>('POST', '/sizing/calculate', request);
		return response.data;
	}

	// Heat Check
	async checkHeat(request: HeatCheckRequest): Promise<HeatCheckResult> {
		const response = await this.request<HeatCheckResult>('POST', '/heat/check', request);
		return response.data;
	}

	// Trade Decisions
	async saveDecision(request: SaveDecisionRequest): Promise<SaveDecisionResponse> {
		const response = await this.request<SaveDecisionResponse>('POST', '/decisions/save', request);
		return response.data;
	}

	// Calendar
	async getCalendar(): Promise<CalendarData> {
		const response = await this.request<CalendarData>('GET', '/calendar');
		return response.data;
	}
}

export interface SizingRequest {
	equity: number;
	risk_pct: number;
	entry: number;
	atr_n: number;
	k: number;
	method: 'stock' | 'opt-delta-atr' | 'opt-maxloss';
	delta?: number;
	max_loss?: number;
}

export interface SizingResult {
	risk_dollars: number;
	stop_distance: number;
	initial_stop: number;
	shares: number;
	contracts: number;
	actual_risk: number;
	method: string;
}

export interface HeatCheckRequest {
	add_risk_dollars: number;
	add_bucket: string;
}

export interface HeatCheckResult {
	current_portfolio_heat: number;
	new_portfolio_heat: number;
	portfolio_heat_pct: number;
	portfolio_cap: number;
	portfolio_cap_exceeded: boolean;
	portfolio_overage: number;
	current_bucket_heat: number;
	new_bucket_heat: number;
	bucket_heat_pct: number;
	bucket_cap: number;
	bucket_cap_exceeded: boolean;
	bucket_overage: number;
	allowed: boolean;
	rejection_reason?: string;
}

export interface SaveDecisionRequest {
	ticker: string;
	entry: number;
	atr: number;
	method: string;
	banner_status: string;
	shares: number;
	contracts: number;
	sector: string;
	strategy: string;
	risk_dollars: number;
	decision: 'GO' | 'NO-GO';
	notes: string;
	banner_green: boolean;
	timer_complete: boolean;
	not_on_cooldown: boolean;
	heat_passed: boolean;
	sizing_complete: boolean;
}

export interface SaveDecisionResponse {
	id: number;
	ticker: string;
	decision: string;
	timestamp: string;
}

export const api = new ApiClient();
