# ffmpeg를 이용한 동영상 편집 서버
동영상 업로드, 트림(Trim), 이어 붙이기(Concat), 정보 조회 기능을 제공합니다.  
  
# Result Code
**200** : 요청 성공  
**400** : 요청 오류  
**500** : 서버 오류  

**1001** : Trim 성공  
**1002** : Concat 성공  
**1003** : Trim 성공 Concat 실패  
**1004** : Concat 성공 Trim 실패  
**1005** : Trim Concat 성공  

**6001** : 허용되지 않는 확장자  
**6002** : Trim 실패  
**6003** : Concat 실패    

**example**  
  
<pre><code>{
    "code": 6002
}</code></pre>

# Script
### `make run`  
9000번 포트로 프로젝트를 실행시킵니다.  
localhost:9000으로 기본 설정이 되어있습니다.  
app_config.yaml 파일에서 포트와 주소를 설정 할 수 있습니다.



# API Example
### POST
### `/upload`  
동영상을 업로드합니다.  
Multipart form 형식으로 업로드합니다.   
app_config.yaml에 설정해놓은 위치로 저장시킵니다.  

**postman example**
  
<img width="834" alt="image" src="https://github.com/SundaePorkCutlet/video-edit/assets/87690981/9ae267ef-1423-4a18-b6d2-a8f09df11319">
  
### POST  
### `/modify`
3가지 조건에 따라 파라미터가 변경됩니다.  


### 1. trim만 요청하는 경우
<pre><code>{
    "isTrimed":true,
    "trimVideoList":[
        {
        "videoId":"2d05c690-ee19-49d2-b256-85a1641cd24b",
        "startTime":2,
        "endTime":5
        },
        {
        "videoId":"210a0287-2d83-4b80-afa5-44cdd709650d",
        "startTime":1,
        "endTime":70
        }
    ]
}</code></pre>

isTrimed : true로 설정
trimVideoList : 여러 동영상 트림 요청 가능 (배열 형식)  
videoId : 요청할 동영상 고유 id  
startTime : 시작 시간 (초 단위)  
endTime : 종료 시간 (초 단위)  

DB trim_history테이블에 생성된 동영상 uuid와 원본 동영상 uuid , startTime, endTime이 저장됩니다.

### 2. concat만 요청하는 경우
<pre><code>{
    "isConcated":true,
    "concatVideoList":[
        "0c649117-6ae9-4d22-b033-80c646191fa5","c801f057-d346-4889-8d5d-68878423f03e"
    ]
}</code></pre>  

  isConcated : true로 설정  
  concatVideoList : 요청할 동영상 고유 id 리스트  

app_config.yaml에 설정해놓은 곳에 원본 동영상 path 리스트가 담긴 txt파일이 생성됩니다.  
DB concat_history테이블에 생성된 동영상 uuid와 encoding된 동영상 uuid 리스트 txt파일 path가 저장됩니다.  

### 3. trim한 동영상들을 concat 요청하는 경우  
<pre><code>{
    "isTrimed":true,
    "trimVideoList":[
         {
        "videoId":"7d04454e-5fd9-48b8-a8bd-ebdabd7e3362",
        "startTime":2,
        "endTime":5
        },
        {
        "videoId":"83158ed9-5ff4-46b2-aed9-63064196b505",
        "startTime":1,
        "endTime":7
        }
    ],
    "isConcated":true
}</code></pre>  

트림과 동일한 요청에 isConcated를 true로 추가  

DB trim_history테이블에 생성된 동영상 uuid와 원본 동영상 uuid , startTime, endTime이 저장됩니다.  
DB encode_history테이블에 encoding된 동영상 uuid와 trim된 원본 동영상 uuid가 저장됩니다.  
DB concat_history테이블에 생성된 동영상 uuid와 encoding된 동영상 uuid 리스트 txt파일 path가 저장됩니다.  
  
### GET
### `/video`  
비디오 관한 모든 정보요청  


#### result example
<pre><code>[
    {
        "video": {
            "id": "105251f7-7516-42fc-80fd-86f0f0ba5368",
            "path": "/Users/hongjunho/Downloads/workspace/stockfolio-test/file/videos/105251f7-7516-42fc-80fd-86f0f0ba5368.mp4",
            "videoName": "trim_KakaoTalk_Video_2024-06-14-20-28-32.mp4",
            "extension": "mp4",
            "uploadTime": "",
            "isTrimed": true,
            "trimTime": "20240617212941",
            "isConcated": false,
            "concatTime": "",
            "isEncoded": false,
            "encodeTime": ""
        },
        "trimInfo": {
            "videoId": "7d04454e-5fd9-48b8-a8bd-ebdabd7e3362",
            "videoPath": "/Users/hongjunho/Downloads/workspace/stockfolio-test/file/videos/7d04454e-5fd9-48b8-a8bd-ebdabd7e3362.mp4",
            "startTime": 2,
            "endTime": 5
        },
        "concatInfoPath": "",
        "encodeInfoPath": ""
    },
]</code></pre>

## ISSUE
동영상을 concat 하기 위해서는 해상도,확장자,코덱을 통일시켜야 했음  
그래서 concat하기전 동영상들을 모두 encoding 시킴  
그래서 trim한 동영상을 concat을 요청을 해버리면 동영상수가 엄청 늘어나는 이슈가 있음  
ex) 3개의 동영상을 trim시킨다음에 concat을 요청함  
1. 3개의 동영상을 trim시킨다음 trim된 동영상들을 저장함 -> 동영상 3개증가, trim_history에 trim된 동영상 uuid와 원본 동영상 uuid를 저장  
2. trim된 3개의 동영상을 concat시키기위해서 encoding시킴 -> 동영상 3개 또 증가, encode_history에 encoding된 동영상 uuid와 trim된 동영상 uuid를 저장  
3. encoding된 동영상 3개를 concat 시킴 -> 동영상 1개 증가, concat txt파일을 생성해서 encoding된 동영상 리스트를 생성, concat_history테이블에 concat된 동영상 uuid와 txt파일 path를 저장  
