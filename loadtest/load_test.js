// k6/load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

const failureRate = new Rate('failed_requests');
const requestDuration = new Trend('request_duration');
const requestsCount = new Counter('total_requests');

export const options = {
    stages: [
        { duration: '5s', target: 50 },
        { duration: '30s', target: 10000 },
        { duration: '10s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<300'],
        failed_requests: ['rate<0.01'],
        http_reqs: ['count>1000'],
    },
};

const BASE_URL = 'http://localhost:8080';

// Тестовые данные
const testTeams = [
    { name: 'backend-team', members: ['user1', 'user2', 'user3', 'user4', 'user5'] },
    { name: 'frontend-team', members: ['user6', 'user7', 'user8', 'user9', 'user10'] },
    { name: 'devops-team', members: ['user11', 'user12', 'user13', 'user14', 'user15'] },
];

const testUsers = [
    { id: 'user1', name: 'Alice' },
    { id: 'user2', name: 'Bob' },
    { id: 'user3', name: 'Charlie' },
    { id: 'user4', name: 'David' },
    { id: 'user5', name: 'Eve' },
    { id: 'user6', name: 'Frank' },
    { id: 'user7', name: 'Grace' },
    { id: 'user8', name: 'Henry' },
    { id: 'user9', name: 'Ivy' },
    { id: 'user10', name: 'Jack' },
    { id: 'user11', name: 'Karen' },
    { id: 'user12', name: 'Leo' },
    { id: 'user13', name: 'Mia' },
    { id: 'user14', name: 'Nathan' },
    { id: 'user15', name: 'Olivia' },
];

// Глобальные переменные для хранения тестовых данных
let testData = {
    teams: [],
    users: [],
    prs: [],
    activeUsers: new Set()
};

export function setup() {
    console.log('Setting up test data...');

    // Создаем команды и пользователей
    testTeams.forEach(team => {
        const members = team.members.map(userId => {
            const user = testUsers.find(u => u.id === userId);
            return {
                user_id: user.id,
                username: user.name,
                is_active: true
            };
        });

        const teamData = {
            team_name: team.name,
            members: members
        };

        const res = http.post(`${BASE_URL}/team/add`, JSON.stringify(teamData), {
            headers: { 'Content-Type': 'application/json' },
        });

        const checkResult = check(res, {
            [`setup: create team ${team.name} status 2xx`]: (r) => r.status >= 200 && r.status < 300,
        });

        if (checkResult && res.status === 201) {
            console.log(`Created team: ${team.name}`);
            testData.teams.push(team.name);
            team.members.forEach(userId => {
                testData.users.push(userId);
                testData.activeUsers.add(userId);
            });
        } else {
            console.log(`Failed to create team ${team.name}: ${res.status} - ${res.body}`);
        }
    });

    sleep(2);

    // Создаем несколько начальных PR
    console.log('Creating initial PRs...');
    for (let i = 0; i < 10; i++) {
        const authorId = testData.users[Math.floor(Math.random() * testData.users.length)];
        const prId = `initial-pr-${i}-${Date.now()}`;

        const prData = {
            pull_request_id: prId,
            pull_request_name: `Initial PR ${i}`,
            author_id: authorId
        };

        const res = http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify(prData), {
            headers: { 'Content-Type': 'application/json' },
        });

        const checkResult = check(res, {
            [`setup: create PR ${prId} status 2xx`]: (r) => r.status >= 200 && r.status < 300,
        });

        if (checkResult && res.status === 201) {
            testData.prs.push({
                id: prId,
                author: authorId,
                status: 'OPEN'
            });
            console.log(`Created PR: ${prId}`);
        } else {
            console.log(`Failed to create PR ${prId}: ${res.status} - ${res.body}`);
        }
        sleep(0.1);
    }

    console.log(`Setup complete: ${testData.teams.length} teams, ${testData.users.length} users, ${testData.prs.length} PRs`);

    return testData;
}

export default function(data) {
    const operations = [
        'get_team',           // 20%
        'create_pr',          // 15%
        'get_user_review',    // 20%
        'merge_pr',           // 10%
        'reassign_reviewer',  // 10%
        'set_user_active',    // 10%
        'get_statistics'      // 15%
    ];

    // Взвешенная вероятность операций
    const weights = [0.20, 0.15, 0.20, 0.10, 0.10, 0.10, 0.15];
    const rand = Math.random();
    let sum = 0;
    let operation = 'get_team';

    for (let i = 0; i < operations.length; i++) {
        sum += weights[i];
        if (rand < sum) {
            operation = operations[i];
            break;
        }
    }

    switch(operation) {
        case 'get_team':
            testGetTeam(data);
            break;
        case 'create_pr':
            testCreatePR(data);
            break;
        case 'get_user_review':
            testGetUserReview(data);
            break;
        case 'merge_pr':
            testMergePR(data);
            break;
        case 'reassign_reviewer':
            testReassignReviewer(data);
            break;
        case 'set_user_active':
            testSetUserActive(data);
            break;
        case 'get_statistics':
            testGetStatistics();
            break;
    }

    sleep(1);
}

function testGetTeam(data) {
    if (data.teams.length === 0) return;

    const teamName = data.teams[Math.floor(Math.random() * data.teams.length)];
    const url = `${BASE_URL}/team/get?team_name=${teamName}`;

    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'get_team' }
    };

    const res = http.get(url, params);

    const checkResult = check(res, {
        'get_team status is 2xx': (r) => r.status >= 200 && r.status < 300,
        'get_team response has team data': (r) => r.json('team_name') !== undefined,
        'get_team response time < 500ms': (r) => r.timings.duration < 500,
    });

    recordMetrics(res, checkResult, 'get_team');
}

function testCreatePR(data) {
    if (data.users.length === 0) return;

    const authorId = data.users[Math.floor(Math.random() * data.users.length)];
    const prId = `pr-${Date.now()}-${Math.random().toString(36).substr(2, 8)}`;

    const prData = {
        pull_request_id: prId,
        pull_request_name: `Load Test PR ${prId}`,
        author_id: authorId
    };

    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'create_pr' }
    };

    const res = http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify(prData), params);

    const checkResult = check(res, {
        'create_pr status is 2xx': (r) => r.status >= 200 && r.status < 300,
        'create_pr response has PR data': (r) => r.json('pr') !== undefined,
        'create_pr response time < 1000ms': (r) => r.timings.duration < 1000,
    });

    if (checkResult && res.status === 201) {
        // Добавляем PR в список для последующих операций
        data.prs.push({
            id: prId,
            author: authorId,
            status: 'OPEN'
        });
    }

    recordMetrics(res, checkResult, 'create_pr');
}

function testGetUserReview(data) {
    if (data.users.length === 0) return;

    const userId = data.users[Math.floor(Math.random() * data.users.length)];
    const url = `${BASE_URL}/users/getReview?user_id=${userId}`;

    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'get_user_review' }
    };

    const res = http.get(url, params);

    const checkResult = check(res, {
        'get_user_review status is 2xx': (r) => r.status >= 200 && r.status < 300,
        'get_user_review response has user data': (r) => r.json('user_id') !== undefined,
        'get_user_review response has PRs array': (r) => Array.isArray(r.json('pull_requests')),
        'get_user_review response time < 500ms': (r) => r.timings.duration < 500,
    });

    recordMetrics(res, checkResult, 'get_user_review');
}

function testMergePR(data) {
    // Ищем OPEN PR для мержа
    const openPRs = data.prs.filter(pr => pr.status === 'OPEN');
    if (openPRs.length === 0) return;

    const pr = openPRs[Math.floor(Math.random() * openPRs.length)];
    const prData = {
        pull_request_id: pr.id
    };

    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'merge_pr' }
    };

    const res = http.post(`${BASE_URL}/pullRequest/merge`, JSON.stringify(prData), params);

    const checkResult = check(res, {
        'merge_pr status is 2xx': (r) => r.status >= 200 && r.status < 300,
        'merge_pr response has PR data': (r) => r.json('pr') !== undefined,
        'merge_pr PR status is MERGED': (r) => r.json('pr.status') === 'MERGED',
        'merge_pr response time < 800ms': (r) => r.timings.duration < 800,
    });

    if (checkResult && res.status === 200) {
        // Обновляем статус PR
        pr.status = 'MERGED';
    }

    recordMetrics(res, checkResult, 'merge_pr');
}

function testReassignReviewer(data) {
    // Ищем OPEN PR для переназначения
    const openPRs = data.prs.filter(pr => pr.status === 'OPEN');
    if (openPRs.length === 0) return;

    const pr = openPRs[Math.floor(Math.random() * openPRs.length)];

    // Получаем информацию о команде чтобы узнать пользователей
    const teamInfoRes = http.get(`${BASE_URL}/team/get?team_name=backend-team`);
    if (teamInfoRes.status !== 200) return;

    const teamInfo = teamInfoRes.json();
    if (!teamInfo || !teamInfo.members || teamInfo.members.length < 2) return;

    // Берем случайного пользователя из команды как "старого ревьювера"
    const oldReviewer = teamInfo.members[Math.floor(Math.random() * teamInfo.members.length)];

    const reassignData = {
        pull_request_id: pr.id,
        old_user_id: oldReviewer.user_id
    };

    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'reassign_reviewer' }
    };

    const res = http.post(`${BASE_URL}/pullRequest/reassign`, JSON.stringify(reassignData), params);

    const checkResult = check(res, {
        'reassign_reviewer status is 2xx': (r) => r.status >= 200 && r.status < 300,
        'reassign_reviewer response has PR data': (r) => r.json('pr') !== undefined,
        'reassign_reviewer response has replaced_by': (r) => r.json('replaced_by') !== undefined,
        'reassign_reviewer response time < 800ms': (r) => r.timings.duration < 800,
    });

    recordMetrics(res, checkResult, 'reassign_reviewer');
}

function testSetUserActive(data) {
    if (data.users.length === 0) return;

    const userId = data.users[Math.floor(Math.random() * data.users.length)];
    // Чередуем активацию/деактивацию
    const isActive = !data.activeUsers.has(userId);

    const userData = {
        user_id: userId,
        is_active: isActive
    };

    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'set_user_active' }
    };

    const res = http.post(`${BASE_URL}/users/setIsActive`, JSON.stringify(userData), params);

    const checkResult = check(res, {
        'set_user_active status is 2xx': (r) => r.status >= 200 && r.status < 300,
        'set_user_active response has user data': (r) => r.json('user') !== undefined,
        'set_user_active user active state is correct': (r) => r.json('user.is_active') === isActive,
        'set_user_active response time < 500ms': (r) => r.timings.duration < 500,
    });

    if (checkResult && res.status === 200) {
        if (isActive) {
            data.activeUsers.add(userId);
        } else {
            data.activeUsers.delete(userId);
        }
    }

    recordMetrics(res, checkResult, 'set_user_active');
}

function testGetStatistics() {
    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'get_statistics' }
    };

    const res = http.get(`${BASE_URL}/statistics/get`, params);

    const checkResult = check(res, {
        'get_statistics status is 2xx': (r) => r.status >= 200 && r.status < 300,
        'get_statistics response has statistics data': (r) => r.json('statistics') !== undefined,
        'get_statistics has total_prs field': (r) => typeof r.json('statistics.total_prs') === 'number',
        'get_statistics has user_assignments array': (r) => Array.isArray(r.json('statistics.user_assignments')),
        'get_statistics response time < 300ms': (r) => r.timings.duration < 300,
    });

    recordMetrics(res, checkResult, 'get_statistics');
}

function recordMetrics(response, checkResult, operation) {
    requestsCount.add(1);
    failureRate.add(!checkResult);
    requestDuration.add(response.timings.duration);

    if (!checkResult) {
        console.log(`Operation ${operation} failed: ${response.status} - ${response.body ? response.body.substring(0, 200) : 'no body'}`);
    }
}

export function teardown(data) {
    console.log('Load test completed');
    console.log(`Final state: ${data.teams.length} teams, ${data.users.length} users, ${data.prs.length} PRs (${data.prs.filter(p => p.status === 'OPEN').length} OPEN, ${data.prs.filter(p => p.status === 'MERGED').length} MERGED)`);

    // Финальные проверки
    const statsRes = http.get(`${BASE_URL}/statistics/get`);
    const finalCheck = check(statsRes, {
        'final statistics status is 2xx': (r) => r.status >= 200 && r.status < 300,
    });

    if (finalCheck) {
        const stats = statsRes.json('statistics');
        console.log(`Final statistics: ${stats.total_prs} total PRs, ${stats.open_prs} open, ${stats.merged_prs} merged`);
    }
}