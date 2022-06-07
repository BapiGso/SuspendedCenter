<?php
set_time_limit(0);
class weblfasrDemo {
    private  $lfasr_host = 'http://raasr.xfyun.cn/api';
    # 请求的接口名
    private $api_prepare = '/prepare';
    private $api_upload = '/upload';
    private $api_merge = '/merge';
    private $api_get_progress = '/getProgress';
    private $api_get_result = '/getResult';
    # 文件分片大小10M
    private $file_piece_sice = 10485760;
    # 转写类型
    private $lfasr_type = 0;
    # 是否开启分词
    private $has_participle = 'false';
    private $has_seperate = 'true';
    # 多候选词个数
    private $max_alternatives = 0;
    # 子用户标识
    private $suid = '';
    private $appid;
    private $secret_key;
    private $upload_file_path;
    private $obj;
    private $data;
    private $sig;
    public function __construct($file_path,$appid,$secret_key,$sig) {
        $this->upload_file_path = $file_path;
        $this->appid = $appid;
        $this->secret_key = $secret_key;
        $this->sig = $sig;
    }

    #主要运行函数
    public function run()
    {
        if ($this->upload_file_path!='') {
            #预处理
            $pre_result = json_decode($this->prepare_request(), true);
            $taskid = $pre_result['data'];
            echo $taskid .'\n';
            #上传文件
            $this->upload_request($taskid, $this->upload_file_path);
            $this->merge_request($taskid);
            while (true) {
                $progress = $this->get_progress_request($taskid);
                $progress_dic = json_decode($progress);
                var_dump($progress_dic);
                if ($progress_dic->err_no != 0 && $progress_dic->err_no != 26605) {
                    return;
                } else {
                    $data = $progress_dic->data;
                    $task_status = json_decode($data);
                    if ($task_status->status == 9) {
                        echo('The task ' . $taskid . ' is in processing, task status: ' . $task_status->desc);
                        echo "\r\n";
                        break;
                    }
                    echo('The task ' . $taskid . ' is in processing, task status: ' . $task_status->desc);
                    echo "\r\n";
                    sleep(1);
                }

            }

            $data_all=$this->get_result_request($taskid);
            $data_all = json_decode($data_all)->data;
            $data_all=json_decode($data_all);
            $message='';
            for ($i=0; $i<count($data_all); $i++) {
                $message.=$data_all[$i]->onebest;
            }
            $this->data=$message;
        }
    }

    public function get_data()
    {
        return $this->data;
    }

    public function gene_params($apiname, $taskid="", $slice_id="")
    {
        $appid = $this->appid;
        $secret_key = $this->secret_key;
        $upload_file_path = $this->upload_file_path;
        $ts=time();
        $md5 = md5($appid.$ts);
        $signa = hash_hmac('sha1',$md5,$secret_key,true);
        $signa=base64_encode($signa);
        $file_len = filesize($this->upload_file_path);
        $file_name = basename($this->upload_file_path);
        if($apiname== $this->api_prepare){
            $slice_num = 1;
            $param_dict['app_id'] = $appid;
            $param_dict['signa'] = $signa;
            $param_dict['ts'] = $ts;
            $param_dict['file_len'] = (string)($file_len);
            $param_dict['file_name'] = $file_name;
            $param_dict['slice_num'] = (string)($slice_num);
        }elseif($apiname== $this->api_upload){
            $param_dict['app_id'] = $appid;
            $param_dict['signa'] = $signa;
            $param_dict['ts'] = $ts;
            $param_dict['task_id'] = $taskid;
            $param_dict['slice_id'] = $slice_id;
        }elseif($apiname == $this->api_merge){
            $param_dict['app_id'] = $appid;
            $param_dict['signa'] = $signa;
            $param_dict['ts'] = $ts;
            $param_dict['task_id'] = $taskid;
            $param_dict['file_name'] = $file_name;
        }elseif($apiname == $this->api_get_progress || $apiname == $this->api_get_result){
            $param_dict['app_id'] = $appid;
            $param_dict['signa'] = $signa;
            $param_dict['ts'] = $ts;
            $param_dict['task_id'] = $taskid;
        }
        return $param_dict;
    }

    public function gene_request($apiname, $data, $files='', $headers=''){
        if($files!=''){
            $data['content'] =$files['content'];
        }
        $jsonStr =$data;
        if($apiname==$this->api_prepare||$apiname==$this->api_merge||$apiname==$this->api_get_progress||$apiname==$this->api_get_result){
            $header= array(
                "Content-Type"=>"application/x-www-form-urlencoded;charset=UTF-8",
            );
        }
        if($apiname==$this->api_upload){
            $header = array(
                "Content-Type"=>"multipart/form-data;",
            );
        }
        $url = $this->lfasr_host . $apiname;
        $rt=$this->http_post_json($url, $jsonStr, $header);
        $result = json_decode($rt[1],true);
        if ($result["ok"] == 0){
            echo("success:".$apiname);
            echo "\r\n";
            return $rt[1];
        }else{
            return $rt[1];
        }
    }
    # 预处理gene_request
    public function prepare_request(){
        return $this->gene_request($apiname=$this->api_prepare,
            $data=$this->gene_params($this->api_prepare));
    }
    #上传
    public function upload_request($taskid, $upload_file_path){
        $finfo    = finfo_open(FILEINFO_MIME_TYPE);
        $mime = finfo_file($finfo, $upload_file_path);
        finfo_close($finfo);
        $curlFile = curl_file_create(
            $upload_file_path,
            $mime,
            pathinfo(
                $upload_file_path,
                PATHINFO_BASENAME
            )
        );
        $files["filename"] = $this->gene_params($this->api_upload)["slice_id"];
        $files["content"] = $curlFile;
        $response = $this->gene_request($this->api_upload,
            $data=$this->gene_params($this->api_upload, $taskid=$taskid,$slice_id=$this->sig->getNextSliceId()),$files=$files);
        return $response;
    }

    public  function merge_request($taskid){
        return $this->gene_request($this->api_merge, $data=$this->gene_params($this->api_merge, $taskid=$taskid));
    }

    # 获取进度
    public  function get_progress_request($taskid){
        return $this->gene_request($this->api_get_progress, $data=$this->gene_params($this->api_get_progress, $taskid=$taskid));
    }

    # 获取结果
    public  function get_result_request($taskid){
        $result=$this->gene_request($this->api_get_result, $data=$this->gene_params($this->api_get_result, $taskid=$taskid));
        return $result;
    }

    #请求
    public  function http_post_json($url, $jsonStr,$header)
    {
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $jsonStr);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $header);
		curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 0);
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        return array($httpCode, $response);
    }
}

class Siggenerator{
    private $__ch;
    public function __construct() {
        $this->__ch = 'aaaaaaaaa`';
    }
    public function getNextSliceId(){
        $ch = $this->__ch;
        $j = strlen($ch) - 1;
        while($j >= 0){
            $cj = $ch[$j];
            if ($cj != 'z'){
                $ch = substr($ch,0,$j) . chr((ord($cj) + 1)) . substr($ch,$j + 1);
                break;
            }
            else{
                $ch = substr($ch,0,$j) . 'a' . substr($ch,$j + 1);
                $j = $j - 1;
            }
        }
        $this->__ch = $ch;
        return $ch;
    }

}

$sig = new Siggenerator();
$app_id = ' ';
$secret_key = ' ';
$file_path = 'resource/0920.wav';
$api = new weblfasrDemo($file_path,$app_id,$secret_key,$sig);
$api->run();
print_r("转写结果：");
print_r($api->get_data());